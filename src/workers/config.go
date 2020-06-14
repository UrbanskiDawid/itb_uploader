package workers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)
//Action that can be done on server
type Action struct {
	Name   string `json:"name"`
	Cmd    string `json:"cmd"`
	Server string `json:"server"`
}

var configurationActions map[string]Action
var configurationServers map[string]Server

//MyConfiguration entire configuration
type MyConfiguration struct {
	Servers []Server `json:"servers"`
	Actions []Action `json:"actions"`
}

//Server definition of remote machine
type Server struct {
	NickName string      `json:"nickname"`
	Host     string      `json:"host"`
	Port     int         `json:"port"`
	Auth     credentials `json:"auth"`
}

type credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func printAction(id int, action *Action) {
	logging.Log.Println("Action #", id)
	logging.Log.Println("Action name: ", action.Name)
	logging.Log.Println("Action server: ", action.Server)
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

func printServer(id int, server *Server) {
	logging.Log.Println("Server #", id)
	logging.Log.Println("Server Nick: ", server.NickName)
	logging.Log.Println("Server Host: ", server.Host)
	logging.Log.Println("Server User: ", server.Auth.User)
	logging.Log.Println("Server Pass: ", server.Auth.Pass)
	logging.Log.Println("Server port: ", server.Port)
}

func LoadConfiguration(cfgFileName string) error {

	//load MyConfiguration from file
	jsonFile, err := os.Open(cfgFileName)
	if err != nil {
		return err
	}

	var cfg MyConfiguration
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return err
	}
	logging.Log.Println("configuration loaded from ", cfgFileName)

	//SERVERS
	configurationServers = make(map[string]Server)

	serversNum := len(cfg.Servers)
	logging.Log.Println("servers found: ", serversNum)
	if serversNum == 0 {
		return errors.New("no servers found in configuration")
	}
	authFromEnviroment := loadEnv()

	for i := 0; i < serversNum; i++ {
		nickName := strings.ToUpper(cfg.Servers[i].NickName)
		overrideServerAuth(&cfg.Servers[i].Auth, &authFromEnviroment)
		configurationServers[nickName] = cfg.Servers[i]
		printServer(i, &cfg.Servers[i])
	}

	//ACTIONS
	configurationActions = make(map[string]Action)

	actionsNum := len(cfg.Actions)
	logging.Log.Println("actions found: ", actionsNum)
	if serversNum == 0 {
		return errors.New("no actions found in configuration")
	}
	for i := 0; i < actionsNum; i++ {
		name := strings.ToUpper(cfg.Actions[i].Name)
		configurationActions[name] = cfg.Actions[i]
		printAction(i, &cfg.Actions[i])
	}

	return nil
}

func getServerByNickName(nickName string) Server {
	nickName = strings.ToUpper(nickName)
	return configurationServers[nickName]
}


func getActionByName(name string) Action {
	name = strings.ToUpper(name)
	return configurationActions[name]
}

func GetActionNames() []string{
    names := make([]string, 0, len(configurationActions))
    for k := range configurationActions {
        names = append(names, k)
    }
	return names
}