package actions

import "github.com/UrbanskiDawid/itb_uploader/logging"

//Action that can be done on server
type Action struct {
	Name         string `json:"name"`
	Cmd          string `json:"cmd"`
	Server       string `json:"server,omitempty"` //omit=localhost
	FileTarget   string `json:"fileTarget,omitempty"`
	FileDownload string `json:"fileDownload,omitempty"`
}

func (a Action) IsLocalAction() bool {
	return a.Server == ""
}

func (a Action) HasUploadFile() bool {
	return a.FileTarget != ""
}

func (a Action) HasDownloadFile() bool {
	return a.FileDownload != ""
}

func (a Action) HasCommand() bool {
	return a.Cmd != ""
}

func printAction(id int, action *Action) {
	logging.Log.Println("Action #", id)
	logging.Log.Println("Action name: ", action.Name)
	logging.Log.Println("Action server: ", action.Server)
	logging.Log.Println("Action target", action.FileTarget)
	logging.Log.Println("Action upload [file]:", action.FileTarget)
	logging.Log.Println("Action download [file]:", action.FileDownload)
}
