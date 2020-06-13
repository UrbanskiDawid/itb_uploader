package main

import "strings"

func executeAction(actionName string) (string, error) {

	action := getActionByName(actionName)

	if strings.ToLower(action.Server) == "localhost" {
		return executeLocal(action.Cmd);
	}

	return executeSSH(action.Cmd, action.Server)
}