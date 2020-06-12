package main

import (
	"fmt"
	"net/http"
)

var cPort int = 8080

func startServer() {
	http.HandleFunc("/", viewIndex)
	http.HandleFunc("/get", viewGet)
	http.HandleFunc("/run", viewRun)
	http.HandleFunc("/ssh", viewSSH)
	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	startServer()
}
