package main

import (
	"fmt"
	"net/http"
	"sync"
)

var html string = `
<html>
<ul>
<li><a href="/get">get</a></li>
<li><a href="/run">run</a></li>
<li><a href="/ssh">ssh</a></li>
<li><a href="/desk/up">upa</a></li>
<li><a href="/desk/down">down</a></li>
</ul>
</html>
`

var servers Servers

type serverStatus struct {
	text string
	num  int
}

var status serverStatus

// ViewIndex main page
func ViewIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html)
	fmt.Fprintf(w, "<p>status: %s</p>", status.text)
	fmt.Fprintf(w, "<p>num: %d</p>", status.num)
	fmt.Fprintf(w, "<p>runCmd: %t %s</p>", dateCmd.running, dateCmd.out)
	fmt.Fprintf(w, "<p>sshCmd: %t %s</p>", voiceCmd.running, voiceCmd.out)
}

// ViewNumber increase number
func ViewNumber(w http.ResponseWriter, r *http.Request) {
	num := status.num
	status.num++

	if num < 10 {
		w.Header().Set("refresh", "1")
	} else {
		w.Header().Set("refresh", "1;url=/")
	}

	fmt.Fprintf(w, "Hello %d ", num)
}

type remoteCmd struct {
	running bool
	lock    sync.Mutex
	out     string
}

var voiceCmd remoteCmd

//ViewVoice play voice via SSH
func ViewVoice(w http.ResponseWriter, r *http.Request) {
	cmd := "python3 /home/dawid/cast.py ASSISTANT_VOICE/dave.mp3"
	serverName := "MINI"

	voiceCmd.lock.Lock()
	voiceCmd.running = true
	go func() {
		defer voiceCmd.lock.Unlock()

		ret, err := executeSSH(cmd, serverName)
		voiceCmd.running = false
		if err == nil {
			voiceCmd.out = ret
		}
	}()

	w.Header().Set("refresh", "1;url=/")
	fmt.Fprint(w, "ssh running")
}

type localCmd struct {
	lock    sync.Mutex
	running bool
	out     string
}

var dateCmd localCmd

//ViewDate get date from remote
func ViewDate(w http.ResponseWriter, r *http.Request) {
	if dateCmd.running {
		fmt.Fprint(w, "busy")
		return
	}
	dateCmd.lock.Lock()
	dateCmd.running = true
	go func() {
		defer dateCmd.lock.Unlock()

		ret, err := executeLocal("date")
		if err == nil {
			dateCmd.out = ret
		}
		dateCmd.running = false
	}()

	w.Header().Set("refresh", "2;url=/")
	fmt.Println("run begin")
	fmt.Fprint(w, "running")
}

var deskCmd remoteCmd

func deskRun(height int) {
	cmd := fmt.Sprintf("sudo /home/dawid/example-moveTo %d", height)
	serverName := "ZERO"

	deskCmd.lock.Lock()
	deskCmd.running = true
	go func() {
		defer deskCmd.lock.Unlock()
		println("cmd:", cmd, serverName)
		ret, err := executeSSH(cmd, serverName)

		println("cmd:", cmd, serverName, " end ")
		deskCmd.running = false
		if err == nil {
			voiceCmd.out = ret
			println("cmd:", cmd, " OK ")
		} else {
			println("cmd:", cmd, " FAIL ", err)
		}
	}()
}

//ViewDeskUp move desk up
func ViewDeskUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("refresh", "3;url=/")
	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRun(6000)
		fmt.Fprint(w, "going up")
	}
}

//ViewDeskDown move desk down
func ViewDeskDown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("refresh", "3;url=/")

	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRun(4000)
		fmt.Fprint(w, "going down")
	}
}
