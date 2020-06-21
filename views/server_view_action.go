package views

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

type viewMemory struct {
	running bool
	lock    sync.Mutex
	out     string
	path    string
}

var actionViewMemory map[string]*viewMemory

//Init must call before using
func Init() {
	actionViewMemory = make(map[string]*viewMemory)
}

const tempFolder string = "TEMP"
const formFileID string = "file"
const maxFileSize int64 = 125000 * 50 //50Mbit

//source: https://medium.com/@petehouston/upload-files-with-curl-93064dcccc76
func SaveFile(r *http.Request) (string, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(maxFileSize)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(formFileID)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//temp folder
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	tempPath := filepath.Join(path, tempFolder)
	os.MkdirAll(tempPath, os.ModePerm)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempPath, "upload-*-"+handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	fullFileName := tempFile.Name() // filepath.Join(".", tempFile.Name())
	return fullFileName, nil
}

func runActionUpload(action base.Action, w http.ResponseWriter, r *http.Request) {

	fileName, err := SaveFile(r)
	if err != nil {
		fmt.Fprint(w, "file error")
	}
	defer os.Remove(fileName)

	logging.Log.Println("uploading file: ", fileName)
	fmt.Println("uploading file: ", fileName)

	err, _ = action.UploadFile(fileName)
	if err != nil {

		logging.Log.Println("fail", err)
		fmt.Println("fail: ", err)
		return
	}
}

func runActionDownload(action base.Action, w http.ResponseWriter, r *http.Request) {

	//temp file
	path, err := os.Getwd()
	if err != nil {
		return
	}
	tempPath := filepath.Join(path, tempFolder)
	os.MkdirAll(tempPath, os.ModePerm)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempPath, "download-*")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(tempFile.Name())

	//remoteFile := executor..GetDownloadFileNameForAction(actionName)
	tmpFilename := tempFile.Name()

	err2, remoteFile := action.DownloadFile(tempFile.Name())
	if err2 != nil {
		logging.Log.Println("fail", err2)
		fmt.Println("fail: ", err2)
		return
	}

	//fmt.Println("Client requests: " + remoteFile)

	//Check if file exists and open
	Openfile, err := os.Open(tmpFilename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(remoteFile))
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client
	return

}

func runActionCommand(action base.Action, w http.ResponseWriter, r *http.Request) string {
	stdOut, stdErr, err := action.Execute()
	if err == nil {
		return stdErr
	}

	logging.Log.Println("runAction cmd: ", action.GetDescription().Name, "end")
	fmt.Println("runAction cmd:", action.GetDescription().Name, "end")
	return stdOut
}

func runAction(action base.Action, w http.ResponseWriter, r *http.Request) {

	mem := actionViewMemory[action.GetDescription().Name]

	if mem.running {
		fmt.Fprint(w, "busy")
	} else {
		mem.lock.Lock()

		mem.running = true

		if action.GetDescription().HasUploadFile() {
			runActionUpload(action, w, r)
		}
		if action.GetDescription().HasDownloadFile() {
			runActionDownload(action, w, r)
		}
		if action.GetDescription().HasCommand() {
			mem.out = runActionCommand(action, w, r)
		}

		fmt.Fprint(w, "srarted")

		defer mem.lock.Unlock()
		mem.running = false
	}

}

//BuildViewAction generate function for server to handle action
func BuildViewAction(action base.Action) func(w http.ResponseWriter, r *http.Request) {

	actionName := action.GetDescription().Name

	actionViewMemory[actionName] = &viewMemory{
		path:    "/action/" + actionName,
		running: false}

	return func(w http.ResponseWriter, r *http.Request) {
		logging.Log.Println("ViewAction", actionName)
		fmt.Println("Request ViewAction", actionName)

		runAction(action, w, r)
	}
}
