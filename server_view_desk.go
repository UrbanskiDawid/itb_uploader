package main

import (
	"fmt"
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

	action := getActionByName(actionName)

	logMsg := fmt.Sprintf("cmd: %s %s", action.Server, action.Cmd)

	deskCmd.lock.Lock()
	deskCmd.running = true
	go func() {
		defer deskCmd.lock.Unlock()
		println(logMsg, "start")
		ret, err := executeSSH(action.Cmd, action.Server)
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
	w.Header().Set("refresh", "3;url=/")

	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRunAction("desk down")
		fmt.Fprint(w, "going down")
	}
}
