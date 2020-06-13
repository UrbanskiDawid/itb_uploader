package workers

import "strings"



func ExecuteAction(actionName string) (string, string, error) {

	action := getActionByName(actionName)

	if strings.ToLower(action.Server) == "localhost" {
		return executeLocal(action.Cmd)
	}

	return executeSSH(action.Cmd, action.Server)
}
