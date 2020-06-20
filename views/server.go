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

	for _, name := range actions.GetActionNames() {

		var actionName string
		actionName = name // note must make a copy

		var userVisibleNameName string
		userVisibleNameName = generateUserVisibleActionName(name)
		http.HandleFunc("/action/"+userVisibleNameName, BuildViewAction(userVisibleNameName, actionName))
	}

	fmt.Println("starting server port", port)
	logging.Log.Println("starting server port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
