package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var configurationServers map[string]Server
var configurationActions map[string]Action

//MyConfiguration entire configuration
type MyConfiguration struct {
	Servers []Server `json:"servers"`
	Actions []Action `json:"actions"`
}

//Action that can be done on server
type Action struct {
	Name   string `json:"name"`
	Cmd    string `json:"cmd"`
	Server string `json:"server"`
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
	fmt.Println("Action #", id)
	fmt.Println("Action name: ", action.Name)
	fmt.Println("Action server: ", action.Server)
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
	fmt.Println("Server #", id)
	fmt.Println("Server Nick: ", server.NickName)
	fmt.Println("Server Host: ", server.Host)
	fmt.Println("Server User: ", server.Auth.User)
	fmt.Println("Server Pass: ", server.Auth.Pass)
	fmt.Println("Server port: ", server.Port)
}

func locadConfiguration(cfgFileName string) error {

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
	fmt.Print("a")

	//SERVERS
	configurationServers = make(map[string]Server)

	serversNum := len(cfg.Servers)
	fmt.Println("servers found: ", serversNum)
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
	fmt.Println("servers found: ", actionsNum)
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
