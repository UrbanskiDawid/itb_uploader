package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type viewDeskMemory struct {
	running bool
	lock    sync.Mutex
	out     string
}

var deskCmd viewDeskMemory

func deskRunAction(actionName string) {

	deskCmd.lock.Lock()
	deskCmd.running = true
	go func() {
		defer deskCmd.lock.Unlock()

		logMsg := fmt.Sprintf("cmd: %s", actionName)
		println(logMsg, "start")
		ret, err := executeAction(actionName)
		deskCmd.running = false
		if err == nil {
			voiceCmd.out = ret
			println(logMsg, " OK ")
		} else {
			println(logMsg, " FAIL ", err)
		}
	}()
}

//ViewDeskUp move desk up
func ViewDeskUp(w http.ResponseWriter, r *http.Request) {
	log.Println("ViewDeskUp")
	w.Header().Set("refresh", "3;url=/")
	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRunAction("desk up")
		fmt.Fprint(w, "going up")
	}
}

//ViewDeskDown move desk down
func ViewDeskDown(w http.ResponseWriter, r *http.Request) {
	log.Println("ViewDeskDown")
	w.Header().Set("refresh", "3;url=/")

	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRunAction("desk down")
		fmt.Fprint(w, "going down")
	}
}
