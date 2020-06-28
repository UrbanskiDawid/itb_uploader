package server

import (
	"fmt"
	"net/http"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

func StartServer(port uint64, act actions.ActionsMap) error {

	actionViewMemory := BuildActionViewMemory()

	http.HandleFunc("/", actionViewMemory.BuildViewIndex())

	actionNames := act.GetNames()

	for _, name := range actionNames {

		var actionName string
		actionName = name // note must make a copy
		action := act.GetByName(name)
		http.HandleFunc("/action/"+actionName, actionViewMemory.BuildViewAction(action))

		fmt.Printf("/action/" + actionName + "\n")
	}

	fmt.Println("starting server port", port)
	logging.Logger.Print("starting server port", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
