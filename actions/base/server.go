package base

import "github.com/UrbanskiDawid/itb_uploader/logging"

//Server definition of remote machine
type Server struct {
	NickName string      `json:"nickname"`
	Host     string      `json:"host"`
	Port     int         `json:"port"`
	Auth     Credentials `json:"auth"`
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func printServer(id int, server *Server) {
	logging.Log.Println("Server #", id)
	logging.Log.Println("Server Nick: ", server.NickName)
	logging.Log.Println("Server Host: ", server.Host)
	logging.Log.Println("Server User: ", server.Auth.User)
	logging.Log.Println("Server Pass: ", server.Auth.Pass)
	logging.Log.Println("Server port: ", server.Port)
}
