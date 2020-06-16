package views

import (
	"fmt"
	"net/http"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

// ViewAllActions show all acions
func ViewAllActions(w http.ResponseWriter, r *http.Request) {
	logging.Log.Println("ViewAllActions")
	fmt.Println("Request ViewAllActions")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, "<h1>actions #%d</h1>", len(actionViewMemory))
	for name := range actionViewMemory {
		fmt.Fprintf(w, "<p><a href=\"%s\">Action %s</a></p>", actionViewMemory[name].path, name)
	}
}
