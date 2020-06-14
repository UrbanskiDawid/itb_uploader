package views

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

type viewMemory struct {
	running bool
	lock    sync.Mutex
	out     string
	path    string
}

var actionViewMemory map[string]viewMemory

//Init must call before using
func Init() {
	actionViewMemory = make(map[string]viewMemory)
}

//BuildViewAction generate function for server to handle action
func BuildViewAction(userVisibleNameName string, actionName string) func(w http.ResponseWriter, r *http.Request) {

	var m viewMemory
	m.path = "/action/" + userVisibleNameName
	m.running = false
	actionViewMemory[actionName] = m

	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Action: %s</br>", actionName)
		logging.Log.Println("ViewAction", actionName)
		fmt.Println("Request ViewAction", actionName)

		w.Header().Set("refresh", "3;url=/")

		mem := actionViewMemory[actionName]

		if mem.running {
			fmt.Fprint(w, "busy")
		} else {
			mem.lock.Lock()

			mem.running = true
			go func() {
				defer mem.lock.Unlock()

				logging.Log.Println("ViewDate cmd: ", actionName, "begin")
				ret, _, err := actions.ExecuteAction(actionName)
				if err == nil {
					mem.out = ret
				}
				mem.running = false
				logging.Log.Println("ViewDate cmd: ", actionName, "end")
				fmt.Println("Request ViewAction", actionName, "end")
			}()

			fmt.Fprint(w, "srarted")
		}
	}
}
