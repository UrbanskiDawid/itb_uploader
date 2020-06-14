package actions

import "strings"

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
