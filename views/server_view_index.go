package views

import (
	"fmt"
	"net/http"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var htmlIndex string = `<html>
<h1>main</h1>
`

func viewError(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotFound)

	logPrefix := fmt.Sprintf("ViewError 404 url: " + r.URL.String())
	logging.LogConsole(logPrefix)
}

// ViewIndex main page
func ViewIndex(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		viewError(w, r)
		return
	}
	url := r.URL.String()

	logPrefix := fmt.Sprintf("Request ViewIndex '%s'", url)
	logging.LogConsole(logPrefix)

	fmt.Fprint(w, htmlIndex)

	fmt.Fprintf(w, "<h1>actions #%d</h1>", len(actionViewMemory))
	for name := range actionViewMemory {
		fmt.Fprintf(w, "<p><a href=\"%s\">Action %s</a></p>", actionViewMemory[name].path, name)
	}
}
