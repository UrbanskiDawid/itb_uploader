package main

import (
	"net/http"
)

var c_port int = 8080

func main() {
	http.HandleFunc("/", view_index)
	http.HandleFunc("/get", view_get)
	http.ListenAndServe(":8080", nil)
}
