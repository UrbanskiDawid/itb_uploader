package main

import (
	"fmt"
	"net/http"
	"os"
)

var cPort int = 8080

func serverInit() {
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	locadConfiguration(jsonFile)
}

func startServer() {

	http.HandleFunc("/", ViewIndex)
	http.HandleFunc("/get", ViewNumber)
	http.HandleFunc("/run", ViewDate)
	http.HandleFunc("/ssh", ViewVoice)
	http.HandleFunc("/desk/up", ViewDeskUp)
	http.HandleFunc("/desk/down", ViewDeskDown)

	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	serverInit()
	startServer()
}
