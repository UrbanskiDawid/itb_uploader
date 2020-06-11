package main

import (
	"fmt"
	"log"
	"net/http"
)

var num int = 0

var html string = `
<html>
<ul>
<li><a href="/get">get</a></li>
<li><a href="/run">run</a></li>
<li><a href="/ssh">ssh</a></li>
</ul>
</html>
`

func viewIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html)
	fmt.Fprintf(w, "status: %s num: %d pid: %d", status.text, status.num, status.pid)
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

func viewSSH(w http.ResponseWriter, r *http.Request) {
	cmd := "python3 /home/dawid/cast.py ASSISTANT_VOICE/dave.mp3"
	remote(cmd)
	w.Header().Set("refresh", "1;url=/")
	fmt.Fprint(w, "ssh running")
}

func viewRun(w http.ResponseWriter, r *http.Request) {
	if status.pid != 0 {
		fmt.Fprint(w, "busy")
		return
	}

	run("date")

	status.text = ""
	status.pid = myCmd.cmd.Process.Pid
	w.Header().Set("refresh", "2;url=/")
	fmt.Fprint(w, "started")

	fmt.Println("run begin")

	// wait for the program to end in a goroutine
	go func() {
		err := myCmd.cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("run end")
		status.text = myCmd.out.String()
		status.pid = 0
		print(status.text)
	}()
}
