package base

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

//Configuration entire configuration
type Configuration struct {
	Servers      []Server      `json:"servers"`
	Descriptions []Description `json:"actions"`
}

func loadEnv() Credentials {
	var ret Credentials
	ret.User = os.Getenv("SSH_USER")
	ret.Pass = os.Getenv("SSH_PASS")
	return ret
}

func overrideServerAuth(auth *Credentials, overrideAuth *Credentials) {
	if overrideAuth.User != "" {
		auth.User = overrideAuth.User
	}
	if overrideAuth.Pass != "" {
		auth.Pass = overrideAuth.Pass
	}
}

func unifyServerName(name string) string {
	name = strings.ToUpper(name)
	if name == "localhost" {
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

func LoadConfigurationFromJson(cfgFileName string) (error, *Configuration) {

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
		cfg.Servers[i].NickName = unifyServerName(cfg.Servers[i].NickName)
	}

	//ACTION DESCRITIONS
	actionsNum := len(cfg.Descriptions)
	logging.Log.Println("action descriptions found: ", actionsNum)
	if serversNum == 0 {
		return errors.New("no actions descriptions found in configuration"), nil
	}
	for i := 0; i < actionsNum; i++ {
		cfg.Descriptions[i].Server = unifyServerName(cfg.Descriptions[i].Server)
		cfg.Descriptions[i].Name = unifyActionName(cfg.Descriptions[i].Name)
	}

	return nil, &cfg
}
