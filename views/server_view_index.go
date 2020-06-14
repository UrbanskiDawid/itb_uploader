package views

import (
	"fmt"
	"net/http"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var htmlIndex string = `<html>
<ul>
<li><a href="/get">get</a></li>
<li><a href="/action/date">date</a></li>
<li><a href="/action/voice">voice</a></li>
</ul>
<p>
<a href="/action/desk/up"><button>desk up</button></a>
<a href="/action/desk/down"><button>desk down</button></a>
</p>
</html>`

// ViewIndex main page
func ViewIndex(w http.ResponseWriter, r *http.Request) {

	logging.Log.Println("ViewIndex")
	fmt.Println("Request ViewIndex")
	
	fmt.Fprint(w, htmlIndex)
	fmt.Fprintf(w, "<p>runNum: %d</p>", viewNumData.num)
	fmt.Fprintf(w, "<p>runCmd: %t %s</p>", dateCmd.running, dateCmd.out)
	fmt.Fprintf(w, "<p>sshCmd: %t %s</p>", voiceCmd.running, voiceCmd.out)
}
