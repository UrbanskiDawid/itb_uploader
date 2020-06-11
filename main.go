package main

import (
	"fmt"
	"net/http"
)

var cPort int = 8080

type Status struct {
	text string
	pid  int
	num  int
}

var status Status

func startServer() {
	http.HandleFunc("/", viewIndex)
	http.HandleFunc("/get", viewGet)
	http.HandleFunc("/run", viewRun)
	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {

	status.text = "start"
	status.pid = 0
	status.num = 0

	startServer()
}
