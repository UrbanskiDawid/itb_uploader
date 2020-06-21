package actions

import (
	"errors"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
)

var ACTIONS ActionsMap

type ActionsMap struct {
	all map[string]base.Action
}

func (e ActionsMap) Init() {
	e.all = make(map[string]base.Action)
}

func (e ActionsMap) GetByName(name string) base.Action {
	return e.all[name]
}

func (e ActionsMap) GetNames() []string {
	keys := make([]string, 0, len(e.all))
	for k := range e.all {
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
	return actionSsh{*action, *server}
}

func buildAllExecutors(
	descriptions []base.Description,
	servers []base.Server,
	onNewExecutor func(base.Action)) error {

	for i := 0; i < len(descriptions); i++ {
		description := &descriptions[i]
		server := findServerIndex(description.Server, servers)
		if server == nil {
			return errors.New("cant find server for:" + description.Name)
		}
		onNewExecutor(actionBuilder(description, server))
	}

	return nil
}

func Init(jsonConfigFile string) error {

	err, cfg := base.LoadConfigurationFromJson(jsonConfigFile)
	if err == nil {
		//ACTIONS.Init()
		ACTIONS.all = make(map[string]base.Action)
		err := buildAllExecutors(cfg.Descriptions, cfg.Servers,
			func(exe base.Action) {
				ACTIONS.all[exe.GetDescription().Name] = exe
			})
		return err
	}
	return err
}
