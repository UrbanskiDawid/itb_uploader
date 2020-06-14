package views

import (
	"fmt"
	"net/http"
	"sync"
	"workers"
	"logging"
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
		logging.Log.Println(logMsg, "start")
		ret, _, err := workers.ExecuteAction(actionName)
		deskCmd.running = false
		if err == nil {
			voiceCmd.out = ret
			logging.Log.Println(logMsg, " OK ")
		} else {
			logging.Log.Println(logMsg, " FAIL ", err)
		}
	}()
}

//ViewDeskUp move desk up
func ViewDeskUp(w http.ResponseWriter, r *http.Request) {
	logging.Log.Println("ViewDeskUp")
	fmt.Println("Request ViewDeskUp")

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
	logging.Log.Println("ViewDeskDown")
	fmt.Println("Request ViewDeskDown")
	
	w.Header().Set("refresh", "3;url=/")

	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRunAction("desk down")
		fmt.Fprint(w, "going down")
	}
}
