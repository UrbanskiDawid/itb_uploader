package views

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

func generateUserVisibleActionName(name string) string {
	var ret string
	ret = name
	ret = strings.ToLower(ret)
	ret = strings.ReplaceAll(ret, " ", "_")
	return ret
}

func StartServer(port uint64) {

	Init()

	http.HandleFunc("/", ViewIndex)
	http.HandleFunc("/action/", ViewIndex)

	actionNames := actions.ACTIONS.GetNames()

	for _, name := range actionNames {

		var actionName string
		actionName = name // note must make a copy
		action := actions.ACTIONS.GetByName(name)
		http.HandleFunc("/action/"+actionName, BuildViewAction(action))
	}

	fmt.Println("starting server port", port)
	logging.Log.Println("starting server port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
