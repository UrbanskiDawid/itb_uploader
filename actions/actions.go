package actions

import (
	"errors"

	"github.com/UrbanskiDawid/itb_uploader/actions/localBackend"
	"github.com/UrbanskiDawid/itb_uploader/actions/sshBackend"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
)

type Action interface {
	Execute() (string, string, error)
	UploadFile(string) error
	DownloadFile(string) error
	GetDescription() base.Description
}

type ActionsMap map[string]Action

func (e ActionsMap) GetByName(name string) Action {
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

func actionBuilder(description *base.Description, server *base.Server) Action {
	if description.IsLocalAction() {
		return localBackend.BuildLocalBackend(*description)
	}

	return sshBackend.BuildActionSsh(*description, *server)
}

func BuildActionMap(descriptions []base.Description, servers []base.Server) (ActionsMap, error) {

	var ACTIONS = make(map[string]Action)

	for i := 0; i < len(descriptions); i++ {
		description := &descriptions[i]
		server := findServerIndex(description.Server, servers)
		if server == nil {
			return ACTIONS, errors.New("cant find server for:" + description.Name)
		}
		exe := actionBuilder(description, server)
		ACTIONS[exe.GetDescription().Name] = exe
	}

	return ACTIONS, nil
}
