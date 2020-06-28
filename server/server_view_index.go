package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var htmlIndex string = ""

//BuildVIeBuildViewIndex show html for all actions
func (actionViewMemory ActionViewMemory) BuildViewIndex(actionNames []string) func(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadFile("../index.html") // just pass the file name
	if err == nil {
		logging.LogConsole("found index.html")

		KNOWN_ACTIONS := "var KNOWN_ACTIONS=["
		for _, name := range actionNames {
			KNOWN_ACTIONS += "'" + name + "',"
		}
		KNOWN_ACTIONS += "];"

		htmlIndex = strings.Replace(string(b), "KNOWN_ACTIONS", KNOWN_ACTIONS, 1)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.String()
		logPrefix := fmt.Sprintf("Request ViewIndex '%s'", url)
		logging.LogConsole(logPrefix)

		fmt.Fprint(w, htmlIndex)

		if htmlIndex == "" {
			fmt.Fprintf(w, "<html><h1>main</h1><h1>actions #%d</h1>", len(actionViewMemory))
			for name := range actionViewMemory {
				fmt.Fprintf(w, "<p><a href=\"%s\">Action %s</a></p>", actionViewMemory[name].path, name)
			}
		}
	}
}
