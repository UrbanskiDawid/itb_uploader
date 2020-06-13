package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Users struct which contains
// an array of users
type Servers struct {
	Server []Server `json:"servers"`
}

// User struct which contains a name
// a type and a list of social links
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

func loadEnv() credentials {
	var ret credentials
	ret.User = os.Getenv("SSH_USER")
	ret.Pass = os.Getenv("SSH_PASS")
	return ret
}

var configuration map[string]Server

func overrideServerAuth(auth *credentials, overrideAuth *credentials) {
	if overrideAuth.User != "" {
		auth.User = overrideAuth.User
	}
	if overrideAuth.Pass != "" {
		auth.Pass = overrideAuth.Pass
	}
}

func printServer(id int, server *Server) {
	fmt.Println("Server #", id, " ", server.NickName)
	fmt.Println("Server Host: ", server.Host)
	fmt.Println("Server User: ", server.Auth.User)
	fmt.Println("Server Pass: ", server.Auth.Pass)
	fmt.Println("Server port: ", server.Port)
}

func locadConfiguration(jsonFile io.Reader) {

	authFromEnviroment := loadEnv()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var servers Servers
	json.Unmarshal(byteValue, &servers)

	configuration = make(map[string]Server)

	for i := 0; i < len(servers.Server); i++ {
		nickName := strings.ToUpper(servers.Server[i].NickName)
		overrideServerAuth(&servers.Server[i].Auth, &authFromEnviroment)
		configuration[nickName] = servers.Server[i]
		printServer(i, &servers.Server[i])
	}
}

func getServerByNickName(nickName string) Server {

	nickName = strings.ToUpper(nickName)
	return configuration[nickName]
}
