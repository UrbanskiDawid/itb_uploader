package main

import (
	"fmt"
	"net/http"
)

var html string = `
<html>
<ul>
<li><a href="/get">get</a></li>
<li><a href="/run">run</a></li>
<li><a href="/ssh">ssh</a></li>
</ul>
</html>
`

type Status struct {
	text string
	num  int
}

var status Status

func viewIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html)
	fmt.Fprintf(w, "<p>status: %s</p>", status.text)
	fmt.Fprintf(w, "<p>num: %d</p>", status.num)
	fmt.Fprintf(w, "<p>runCmd: %t %s</p>", runCmd.running, runCmd.out)
	fmt.Fprintf(w, "<p>sshCmd: %t %s</p>", remoteCmd.running, remoteCmd.out)
}

func viewGet(w http.ResponseWriter, r *http.Request) {
	num := status.num
	status.num++

	if num < 10 {
		w.Header().Set("refresh", "1")
	} else {
		w.Header().Set("refresh", "1;url=/")
	}

	fmt.Fprintf(w, "Hello %d ", num)
}

type RemoteCmd struct {
	running bool
	out     string
}

var remoteCmd RemoteCmd

func viewSSH(w http.ResponseWriter, r *http.Request) {
	cmd := "python3 /home/dawid/cast.py ASSISTANT_VOICE/dave.mp3"

	remoteCmd.running = true
	go func() {
		ret, err := remote(cmd)
		remoteCmd.running = false
		if err == nil {
			remoteCmd.out = ret
		}
	}()

	w.Header().Set("refresh", "1;url=/")
	fmt.Fprint(w, "ssh running")
}

type MyCmd struct {
	running bool
	out     string
}

var runCmd MyCmd

func viewRun(w http.ResponseWriter, r *http.Request) {
	if runCmd.running {
		fmt.Fprint(w, "busy")
		return
	}

	runCmd.running = true
	go func() {
		ret, err := run("date")
		if err == nil {
			runCmd.out = ret
		}
		runCmd.running = false
	}()

	w.Header().Set("refresh", "2;url=/")
	fmt.Println("run begin")
	fmt.Fprint(w, "running")
}
