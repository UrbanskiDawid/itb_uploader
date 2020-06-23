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

func openTmpFile(postfix string) (*os.File, error) {
	//temp folder
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	tempPath := filepath.Join(path, tempFolder)
	os.MkdirAll(tempPath, os.ModePerm)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempPath, "upload-*-"+postfix) //handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	return tempFile, nil
}

func runActionUpload(action base.Action, w http.ResponseWriter, r *http.Request) string {

	logPrefix := fmt.Sprintf("runActionUpload() [%s@%s] ", action.GetDescription().FileTarget, action.GetDescription().Server)
	logging.LogConsole(logPrefix + "BEGIN")

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		logging.LogConsole(logPrefix + fmt.Sprintf("erorr: missing file"))
		return "error: missing file"
	}
	receivedFile, handler, err := r.FormFile(formFileID)
	if err != nil {
		logging.LogConsole(logPrefix + fmt.Sprintf("error: file read"))
		return "error: file read"
	}
	defer receivedFile.Close()

	logging.LogConsole(logPrefix + fmt.Sprintf("received: '%s' uploading...", handler.Filename))

	err2, remoteFile := action.UploadFile(receivedFile)
	if err2 != nil {
		logging.LogConsole(logPrefix + fmt.Sprint("upload file failed", err))
		return "error: upload failed"
	}

	logging.LogConsole(logPrefix + "END")
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

	err2, remoteFile := action.DownloadFile(tempFile)
	if err2 != nil {
		logging.LogConsole(logMsg + fmt.Sprint("can't download from remote:", err2))
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
