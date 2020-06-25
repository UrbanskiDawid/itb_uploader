package base

import (
	"fmt"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

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
	logging.LogConsole(fmt.Sprintf("Action #%d", id))
	logging.LogConsole(fmt.Sprintf("Action name: %s", action.Name))
	logging.LogConsole(fmt.Sprintf("Action server: %s", action.Server))
	logging.LogConsole(fmt.Sprintf("Action target %s", action.FileTarget))
	logging.LogConsole(fmt.Sprintf("Action upload [file]: %s", action.FileTarget))
	logging.LogConsole(fmt.Sprintf("Action download [file]: %s", action.FileDownload))
}
