package main

import (
	"fmt"
	"net/http"
)

var cPort int = 8080

func serverInit() {
	fmt.Println("starting Init")
	err := locadConfiguration("config.json")
	if err != nil {
		panic(err)
	}
}

var htmlIndex string = `
<html>
<ul>
<li><a href="/get">get</a></li>
<li><a href="/action/date">date</a></li>
<li><a href="/action/voice">voice</a></li>
</ul>
<p>
<a href="/action/desk/up"><button>desk up</button></a>
<a href="/action/desk/down"><button>desk donw</button></a>
</p>
</html>
`

func startServer() {

	http.HandleFunc("/", ViewIndex)
	http.HandleFunc("/get", ViewNumber)
	http.HandleFunc("/action/date", ViewDate)
	http.HandleFunc("/action/voice", ViewVoice)
	http.HandleFunc("/action/desk/up", ViewDeskUp)
	http.HandleFunc("/action/desk/down", ViewDeskDown)

	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	serverInit()
	startServer()
}
