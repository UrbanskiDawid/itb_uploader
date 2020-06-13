package main

import (
	"fmt"
	"net/http"
)

var cPort int = 8080

func serverInit() {
	fmt.Println("starting Init")

	configFileName, err := FindConfigFileName()
	if err != nil {
		panic(err)
	}

	err = loadConfiguration(configFileName)
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
	http.HandleFunc("/upload", ViewUploadFile)

	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	serverInit()
	startServer()
}
