package views

import (
	"fmt"
	"net/http"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var htmlIndex string = `<html>
<h1>main</h1>
<ul>
<li><a href="/action/">actions</a></li>
</ul>
</html>`

// ViewIndex main page
func ViewIndex(w http.ResponseWriter, r *http.Request) {

	logging.Log.Println("ViewIndex")
	fmt.Println("Request ViewIndex")

	fmt.Fprint(w, htmlIndex)

	fmt.Fprintf(w, "<h1>actions #%d</h1>", len(actionViewMemory))
	for name := range actionViewMemory {
		fmt.Fprintf(w, "<p><a href=\"%s\">Action %s</a></p>", actionViewMemory[name].path, name)
	}
}
