package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

//MyConfiguration entire configuration
type MyConfiguration struct {
	Servers []Server `json:"servers"`
	Actions []Action `json:"actions"`
}

func loadEnv() credentials {
	var ret credentials
	ret.User = os.Getenv("SSH_USER")
	ret.Pass = os.Getenv("SSH_PASS")
	return ret
}

func overrideServerAuth(auth *credentials, overrideAuth *credentials) {
	if overrideAuth.User != "" {
		auth.User = overrideAuth.User
	}
	if overrideAuth.Pass != "" {
		auth.Pass = overrideAuth.Pass
	}
}

func overrideServerName(name string) string {
	name = strings.ToUpper(name)
	if name == "localhost" {
		name = ""
	}
	return name
}

func loadConfigurationFromJson(cfgFileName string) (error, *MyConfiguration) {

	//load MyConfiguration from file
	jsonFile, err := os.Open(cfgFileName)
	if err != nil {
		return err, nil
	}

	var cfg MyConfiguration
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return err, nil
	}
	logging.Log.Println("configuration loaded from ", cfgFileName)

	//SERVERS
	serversNum := len(cfg.Servers)
	logging.Log.Println("servers found: ", serversNum)
	if serversNum == 0 {
		return errors.New("no servers found in configuration"), nil
	}
	authFromEnviroment := loadEnv()

	for i := 0; i < serversNum; i++ {
		overrideServerAuth(&cfg.Servers[i].Auth, &authFromEnviroment)
		cfg.Servers[i].NickName = overrideServerName(cfg.Servers[i].NickName)
	}

	//ACTIONS
	actionsNum := len(cfg.Actions)
	logging.Log.Println("actions found: ", actionsNum)
	if serversNum == 0 {
		return errors.New("no actions found in configuration"), nil
	}
	for i := 0; i < actionsNum; i++ {
		cfg.Actions[i].Server = overrideServerName(cfg.Actions[i].Server)
	}

	return nil, &cfg
}
