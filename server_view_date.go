package main

import (
	"fmt"
	"net/http"
	"sync"
)

type viewPackageMemory struct {
	lock    sync.Mutex
	running bool
	out     string
}

var dateCmd viewPackageMemory

//ViewDate get date from remote
func ViewDate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("refresh", "2;url=/")

	if dateCmd.running {
		fmt.Fprint(w, "busy")
		return
	}

	action := getActionByName("date")

	dateCmd.lock.Lock()
	dateCmd.running = true
	go func() {
		defer dateCmd.lock.Unlock()

		fmt.Println("ViewDate cmd: ", action.Cmd, "begin")
		ret, err := executeLocal(action.Cmd)
		if err == nil {
			dateCmd.out = ret
		}
		dateCmd.running = false
		fmt.Println("ViewDate cmd: ", action.Cmd, "end")
	}()

	fmt.Fprint(w, "running")
}
