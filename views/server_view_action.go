package views

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/UrbanskiDawid/itb_uploader/tmp"
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

const formFileID string = "file"
const maxFileSize int64 = 32 << 10

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

	//create tmp file
	tmpFile, err := tmp.OpenTmpFile("upload" + handler.Filename)
	if err != nil {
		logging.LogConsole(logPrefix + fmt.Sprintf("error: cant create tmp file"))
		return "error: file read"
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	//LIFO

	err2 := action.UploadFile(tmpFile.Name())
	if err2 != nil {
		logging.LogConsole(logPrefix + fmt.Sprint("upload file failed", err))
		return "error: upload failed"
	}

	logging.LogConsole(logPrefix + "END")
	return action.GetDescription().FileTarget
}

func runActionDownload(action base.Action, w http.ResponseWriter, r *http.Request) {

	logMsg := fmt.Sprintf("runActionDownload() ['%s' from '%s'] ", action.GetDescription().FileDownload, action.GetDescription().Server)

	logging.LogConsole(logMsg + "begin")

	remoteBaseFileName := path.Base(action.GetDescription().FileDownload)

	tmpFileName := tmp.GenerateTmpFileName(remoteBaseFileName)
	if tmpFileName == "" {
		logging.LogConsole(logMsg + "failed to generate tmp file name")
		return
	}

	//============================= GET FILE ========================
	err := action.DownloadFile(tmpFileName)
	if err != nil {
		logging.LogConsole(logMsg + fmt.Sprint("can't download from remote:", err))
		return
	}
	logging.LogConsole(logMsg + "download done, sending back")
	defer os.Remove(tmpFileName)

	//============================ SEND FILE =======================
	tmpFile, err := os.Open(tmpFileName)
	if err != nil {
		logging.LogConsole(logMsg + fmt.Sprintf("local file not found: %s", tmpFileName))
		return
	}
	defer tmpFile.Close()

	//============= READ
	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	tmpFile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := tmpFile.Stat()                      //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(action.GetDescription().FileDownload))
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	tmpFile.Seek(0, 0)
	io.Copy(w, tmpFile) //'Copy' the file to the client
	//==============================================================

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
