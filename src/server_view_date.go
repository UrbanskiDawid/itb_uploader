package main

import (
	"fmt"
	"net/http"
	"sync"
	workers "dawidurbanski.pl/itb_uploader/workers"
)

type viewPackageMemory struct {
	lock    sync.Mutex
	running bool
	out     string
}

var dateCmd viewPackageMemory

//ViewDate get date from remote
func ViewDate(w http.ResponseWriter, r *http.Request) {

	Log.Println("ViewDate")

	actionName := "date"
	w.Header().Set("refresh", "2;url=/")

	if dateCmd.running {
		fmt.Fprint(w, "busy")
		return
	}

	dateCmd.lock.Lock()
	dateCmd.running = true
	go func() {
		defer dateCmd.lock.Unlock()

		fmt.Println("ViewDate cmd: ", actionName, "begin")
		ret, _, err := workers.ExecuteAction(actionName)
		if err == nil {
			dateCmd.out = ret
		}
		dateCmd.running = false
		fmt.Println("ViewDate cmd: ", actionName, "end")
	}()

	fmt.Fprint(w, "running")
}
