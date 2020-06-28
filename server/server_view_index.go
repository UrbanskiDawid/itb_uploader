package server

import (
	"fmt"
	"net/http"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var htmlIndex string = `<html>
<h1>main</h1>
`

//BuildVIeBuildViewIndex show html for all actions
func (actionViewMemory ActionViewMemory) BuildViewIndex() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.String()
		logPrefix := fmt.Sprintf("Request ViewIndex '%s'", url)
		logging.LogConsole(logPrefix)

		fmt.Fprint(w, htmlIndex)

		fmt.Fprintf(w, "<h1>actions #%d</h1>", len(actionViewMemory))
		for name := range actionViewMemory {
			fmt.Fprintf(w, "<p><a href=\"%s\">Action %s</a></p>", actionViewMemory[name].path, name)
		}
	}
}
