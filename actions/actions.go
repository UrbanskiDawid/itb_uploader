package actions

import "strings"

var configurationActions map[string]*Action
var configurationServers map[string]Server

func getServerByNickName(nickName string) Server {
	nickName = strings.ToUpper(nickName)
	return configurationServers[nickName]
}

func getActionByName(name string) *Action {
	name = strings.ToUpper(name)
	return configurationActions[name]
}

func GetTargetFileNameForAction(actionName string) string {
	return configurationActions[actionName].FileTarget
}

func GetSourceFileNameForAction(actionName string) string {
	return configurationActions[actionName].FileDownload
}

func IsActionWithUploadFile(name string) bool {
	return configurationActions[name].FileTarget != ""
}

func IsActionWithDownloadFile(name string) bool {
	return configurationActions[name].FileDownload != ""
}

func isLocalhost(action *Action) bool {
	return strings.ToLower(action.Server) == "localhost"
}

func ExecuteAction(actionName string) (string, string, error) {

	action := getActionByName(actionName)
	if isLocalhost(action) {
		return executeLocal(action.Cmd)
	}
	return executeSSH(action.Cmd, action.Server)
}

func GetActionNames() []string {
	names := make([]string, 0, len(configurationActions))
	for k := range configurationActions {
		names = append(names, k)
	}
	return names
}

func UploadFile(actionName string, localFile string) error {
	action := getActionByName(actionName)
	if isLocalhost(action) {
		return copyFileLocal(localFile, action.FileTarget)
	}

	return uploadFileSSH(action.Server, localFile, action.FileTarget)
}

func DownloadFile(actionName string, localFile string) error {
	action := getActionByName(actionName)
	if isLocalhost(action) {
		return copyFileLocal(action.FileDownload, localFile)
	}
	return downloadFileSSH(action.Server, localFile, action.FileDownload)
}
