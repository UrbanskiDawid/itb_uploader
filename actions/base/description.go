package base

import "github.com/UrbanskiDawid/itb_uploader/logging"

//ActionConfiguration that can be done on server
type Description struct {
	Name         string `json:"name"`
	Cmd          string `json:"cmd"`
	Server       string `json:"server,omitempty"` //omit=localhost
	FileTarget   string `json:"fileTarget,omitempty"`
	FileDownload string `json:"fileDownload,omitempty"`
}

func (a Description) IsLocalAction() bool {
	return a.Server == ""
}

func (a Description) HasUploadFile() bool {
	return a.FileTarget != ""
}

func (a Description) HasDownloadFile() bool {
	return a.FileDownload != ""
}

func (a Description) HasCommand() bool {
	return a.Cmd != ""
}

func printAction(id int, action *Description) {
	logging.Log.Println("Action #", id)
	logging.Log.Println("Action name: ", action.Name)
	logging.Log.Println("Action server: ", action.Server)
	logging.Log.Println("Action target", action.FileTarget)
	logging.Log.Println("Action upload [file]:", action.FileTarget)
	logging.Log.Println("Action download [file]:", action.FileDownload)
}
