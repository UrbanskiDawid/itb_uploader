package main

import (
	"fmt"
	"net/http"
)

var cPort int = 8080

func serverInit() {
	fmt.Println("starting Init")
	err := locadConfiguration("config.json")
	if err != nil {
		panic(err)
	}
}

func startServer() {

	http.HandleFunc("/", ViewIndex)
	http.HandleFunc("/get", ViewNumber)
	http.HandleFunc("/action/date", ViewDate)
	http.HandleFunc("/action/voice", ViewVoice)
	http.HandleFunc("/action/desk/up", ViewDeskUp)
	http.HandleFunc("/action/desk/down", ViewDeskDown)

	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	serverInit()
	startServer()
}
