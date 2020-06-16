package actions

import "strings"

var configurationActions map[string]Action
var configurationServers map[string]Server

func getServerByNickName(nickName string) Server {
	nickName = strings.ToUpper(nickName)
	return configurationServers[nickName]
}

func getActionByName(name string) Action {
	name = strings.ToUpper(name)
	return configurationActions[name]
}

func ExecuteAction(actionName string) (string, string, error) {

	action := getActionByName(actionName)

	if strings.ToLower(action.Server) == "localhost" {
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
