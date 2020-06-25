package base

import (
	"fmt"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

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
	logging.LogConsole(fmt.Sprintf("Server #%d", id))
	logging.LogConsole("Server Nick: " + server.NickName)
	logging.LogConsole("Server Host: " + server.Host)
	logging.LogConsole("Server User: " + server.Auth.User)
	logging.LogConsole("Server Pass: " + server.Auth.Pass)
	logging.LogConsole(fmt.Sprintf("Server port: %d", server.Port))
}
