package views

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
const maxFileSize int64 = 32 << 10

//source: https://medium.com/@petehouston/upload-files-with-curl-93064dcccc76
func SaveFile(r *http.Request) (string, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		return "", err
	}
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

func runActionUpload(action base.Action, w http.ResponseWriter, r *http.Request) string {

	logging.LogConsole(fmt.Sprint("runActionUpload() file: ", action.GetDescription().FileTarget, " to "+action.GetDescription().Server))

	fileName, err := SaveFile(r)
	if err != nil {
		logging.LogConsole(fmt.Sprintf("upload failed to recive file"))
		fmt.Fprint(w, "please upload file")
		return ""
	}
	defer os.Remove(fileName)

	logging.LogConsole(fmt.Sprint("uploading file: ", action.GetDescription().FileTarget, " to "+action.GetDescription().Server))

	err2, remoteFile := action.UploadFile(fileName)
	if err2 != nil {

		logging.LogConsole(fmt.Sprint("upload file failed", err))
		return "fail"
	}

	logging.LogConsole(fmt.Sprint("uploading file: ", action.GetDescription().FileTarget, " to "+action.GetDescription().Server, "OK"))
	return remoteFile
}

func runActionDownload(action base.Action, w http.ResponseWriter, r *http.Request) {

	logMsg := fmt.Sprintf("runActionDownload() ['%s' from '%s'] ", action.GetDescription().FileDownload, action.GetDescription().Server)

	logging.LogConsole(logMsg + "begin")

	//temp file
	path, err := os.Getwd()
	if err != nil {
		logging.LogConsole(logMsg + "failed to get working dir")
		return
	}
	tempPath := filepath.Join(path, tempFolder)
	os.MkdirAll(tempPath, os.ModePerm)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempPath, "download-*")
	if err != nil {
		logging.LogConsole(logMsg + fmt.Sprint("can't create local temp file:", err))
		return
	}
	defer os.Remove(tempFile.Name())
	logging.LogConsole(logMsg + fmt.Sprintf("local temp file %s created", tempFile.Name()))

	//remoteFile := executor..GetDownloadFileNameForAction(actionName)
	tmpFilename := tempFile.Name()

	err2, remoteFile := action.DownloadFile(tempFile.Name())
	if err2 != nil {
		logging.LogConsole(logMsg + fmt.Sprint("cant download from remote:", err2))
		return
	}

	logging.LogConsole(logMsg + "download done, sending back")

	//fmt.Println("Client requests: " + remoteFile)

	//Check if file exists and open
	Openfile, err := os.Open(tmpFilename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		logging.LogConsole(logMsg + fmt.Sprintf("local file not found: %s", tmpFilename))

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

	logging.LogConsole(logMsg + "sending done")
	return
}

func runActionCommand(action base.Action, w http.ResponseWriter, r *http.Request) string {

	logPrefix := fmt.Sprintf("runActionCommand() action:%s cmd:%s ", action.GetDescription().Name, action.GetDescription().Cmd)

	logging.LogConsole(logPrefix + "BEGIN")

	stdOut, stdErr, err := action.Execute()

	logging.LogConsole(logPrefix + fmt.Sprintf("stdOut:%s end stdErr:%s", strings.ReplaceAll(stdOut, "\n", "\\n"), strings.ReplaceAll(stdErr, "\n", "\\n")))

	if err == nil {
		return stdErr
	}
	return stdOut
}

func runAction(action base.Action, w http.ResponseWriter, r *http.Request) {

	mem := actionViewMemory[action.GetDescription().Name]

	if mem.running {
		fmt.Fprint(w, "busy")
	} else {
		mem.lock.Lock()

		mem.running = true

		defer mem.lock.Unlock()

		if action.GetDescription().HasUploadFile() {
			mem.out = runActionUpload(action, w, r)
			fmt.Fprint(w, mem.out)
		}
		if action.GetDescription().HasDownloadFile() {
			runActionDownload(action, w, r)
		}
		if action.GetDescription().HasCommand() {
			mem.out = runActionCommand(action, w, r)
			fmt.Fprint(w, mem.out)
		}

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

		logging.LogConsole("ViewAction" + actionName + " BEGIN")
		runAction(action, w, r)
		logging.LogConsole("ViewAction" + actionName + " END")
	}
}
