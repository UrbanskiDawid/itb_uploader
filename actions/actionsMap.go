package actions

import (
	"errors"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
)

type ActionsMap map[string]base.Action

func BuildActionMap() ActionsMap {
	var ACTIONS = make(map[string]base.Action)
	return ACTIONS
}

func (e ActionsMap) GetByName(name string) base.Action {
	return e[name]
}

func (e ActionsMap) GetNames() []string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	return keys
}

func findServerIndex(name string, servers []base.Server) *base.Server {

	for j := 0; j < len(servers); j++ {
		if servers[j].NickName == name {
			return &servers[j]
		}
	}
	return nil
}

func actionBuilder(action *base.Description, server *base.Server) base.Action {
	if action.IsLocalAction() {
		return ActionLocal{*action}
	}

	client, err := buildClientConfig(*server)
	if err != nil {
		panic("server " + server.NickName + " configuration error")
	}
	return actionSsh{*action, *server, *client}
}

func (e ActionsMap) BuildAllExecutors(descriptions []base.Description, servers []base.Server) error {

	for i := 0; i < len(descriptions); i++ {
		description := &descriptions[i]
		server := findServerIndex(description.Server, servers)
		if server == nil {
			return errors.New("cant find server for:" + description.Name)
		}
		exe := actionBuilder(description, server)
		e[exe.GetDescription().Name] = exe
	}

	return nil
}
