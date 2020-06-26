package main

import (
	"encoding/json"
	"errors"
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

	return nil, &cfg
}

func fixServersInConfiguration(cfg *Configuration) {

	var localhostServer base.Server
	localhostServer.NickName = ""
	localhostServer.Port = 0
	cfg.Servers = append(cfg.Servers, localhostServer)

	serversNum := len(cfg.Servers)
	for i := 0; i < serversNum; i++ {
		cfg.Servers[i].NickName = unifyServerName(cfg.Servers[i].NickName)
	}
}

func fixActionsInConfiguration(cfg *Configuration) {
	actionsNum := len(cfg.Descriptions)
	for i := 0; i < actionsNum; i++ {
		cfg.Descriptions[i].Server = unifyServerName(cfg.Descriptions[i].Server)
		cfg.Descriptions[i].Name = unifyActionName(cfg.Descriptions[i].Name)
	}
}

//InitConfig load configuration form json file
func InitConfig(jsonConfigFile string) (actions.ActionsMap, error) {

	err, cfg := loadConfigurationFromJson(jsonConfigFile)
	if err != nil {
		return nil, err
	}

	if len(cfg.Servers) == 0 {
		return nil, errors.New("no servers found in configuration")
	}

	if len(cfg.Descriptions) == 0 {
		return nil, errors.New("no actions descriptions found in configuration")
	}

	fixServersInConfiguration(cfg)
	fixActionsInConfiguration(cfg)

	return actions.BuildActionMap(cfg.Descriptions, cfg.Servers)
}
