package views

import (
	"fmt"
	"net/http"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var htmlIndex string = `<html>
<h1>main</h1>
<ul>
<li><a href="/get">get</a></li>
<li><a href="/action/">actions</a></li>
</ul>
</html>`

// ViewIndex main page
func ViewIndex(w http.ResponseWriter, r *http.Request) {

	logging.Log.Println("ViewIndex")
	fmt.Println("Request ViewIndex")

	fmt.Fprint(w, htmlIndex)
}
