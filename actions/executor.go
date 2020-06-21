package actions

import (
	"errors"
	"strings"
)

type Executor interface {
	Execute() (string, string, error)
	UploadFile(localFile string) (error, string)
	DownloadFile(localFile string) (error, string)
	GetAction() Action
}

func executorFactory(action *Action, server *Server) Executor {
	if action.IsLocalAction() {
		return ExecutorLocal{*action}
	}
	return executorSsh{*action, *server}
}

func findServerIndex(name string, servers []Server) *Server {

	for j := 0; j < len(servers); j++ {
		if servers[j].NickName == name {
			return &servers[j]
		}
	}
	return nil
}

func buildAllExecutors(
	actions []Action,
	servers []Server,
	onNewExecutor func(string, Executor)) error {

	for i := 0; i < len(actions); i++ {
		action := &actions[i]
		actionName := strings.ToLower(action.Name)
		server := findServerIndex(action.Server, servers)
		if server == nil {
			return errors.New("cant find server for:" + actionName)
		}
		onNewExecutor(actionName, executorFactory(action, server))
	}

	return nil
}
