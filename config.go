package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

//Configuration entire configuration
type Configuration struct {
	Servers      []base.Server      `json:"servers"`
	Descriptions []base.Description `json:"actions"`
}

func overrideServerAuth(auth *base.Credentials, overrideAuth *base.Credentials) {
	if overrideAuth.User != "" {
		auth.User = overrideAuth.User
	}
	if overrideAuth.Pass != "" {
		auth.Pass = overrideAuth.Pass
	}
}

func unifyServerName(name string) string {
	name = strings.ToUpper(name)
	if name == "LOCALHOST" {
		name = ""
	}
	return name
}

var re = regexp.MustCompile(`[^0-9A-Za-z_]`)

func unifyActionName(name string) string {
	name = strings.ToLower(name)
	name = re.ReplaceAllString(name, "_")
	return name
}

func loadConfigurationFromJson(cfgFileName string) (error, *Configuration) {

	//load MyConfiguration from file
	jsonFile, err := os.Open(cfgFileName)
	if err != nil {
		return err, nil
	}

	var cfg Configuration
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return err, nil
	}
	logging.LogConsole("configuration loaded from " + cfgFileName)

	serversNum := len(cfg.Servers)
	logging.LogConsole(fmt.Sprintf("servers found: %d", serversNum))
	if serversNum == 0 {
		return errors.New("no servers found in configuration"), nil
	}

	actionsNum := len(cfg.Descriptions)
	logging.LogConsole(fmt.Sprintf("action descriptions found:%d", actionsNum))
	if actionsNum == 0 {
		return errors.New("no actions descriptions found in configuration"), nil
	}

	return nil, &cfg
}

func fixServersInConfiguration(cfg *Configuration) {

	serversNum := len(cfg.Servers)

	//Add localhost
	var localhostServer base.Server
	localhostServer.NickName = ""
	localhostServer.Port = 0
	cfg.Servers = append(cfg.Servers, localhostServer)

	for i := 0; i < serversNum; i++ {
		cfg.Servers[i].NickName = unifyServerName(cfg.Servers[i].NickName)
	}
}

func fixActionsInConfiguration(cfg *Configuration) {
	actionsNum := len(cfg.Descriptions)
	logging.LogConsole(fmt.Sprintf("action descriptions found:%d", actionsNum))

	for i := 0; i < actionsNum; i++ {
		cfg.Descriptions[i].Server = unifyServerName(cfg.Descriptions[i].Server)
		cfg.Descriptions[i].Name = unifyActionName(cfg.Descriptions[i].Name)
	}
}

func InitConfig(jsonConfigFile string) (actions.ActionsMap, error) {

	ACTIONS := actions.BuildActionMap()

	err, cfg := loadConfigurationFromJson(jsonConfigFile)
	if err != nil {
		return ACTIONS, err
	}

	fixServersInConfiguration(cfg)
	fixActionsInConfiguration(cfg)

	ACTIONS.BuildAllExecutors(cfg.Descriptions, cfg.Servers)
	return ACTIONS, err
}
