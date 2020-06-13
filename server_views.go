package main

import (
	"fmt"
	"net/http"
	"sync"
)

// ViewIndex main page
func ViewIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlIndex)
	fmt.Fprintf(w, "<p>runNum: %d</p>", viewNumData.num)
	fmt.Fprintf(w, "<p>runCmd: %t %s</p>", dateCmd.running, dateCmd.out)
	fmt.Fprintf(w, "<p>sshCmd: %t %s</p>", voiceCmd.running, voiceCmd.out)
}

type viewNumberData struct {
	num int
}

var viewNumData viewNumberData

// ViewNumber increase number
func ViewNumber(w http.ResponseWriter, r *http.Request) {
	num := viewNumData.num

	if num < 5 {
		viewNumData.num++
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

	w.Header().Set("refresh", "2;url=/")

	if dateCmd.running {
		fmt.Fprint(w, "busy")
		return
	}

	action := getActionByName("date")

	dateCmd.lock.Lock()
	dateCmd.running = true
	go func() {
		defer dateCmd.lock.Unlock()

		fmt.Println("ViewDate cmd: ", action.Cmd, "begin")
		ret, err := executeLocal(action.Cmd)
		if err == nil {
			dateCmd.out = ret
		}
		dateCmd.running = false
		fmt.Println("ViewDate cmd: ", action.Cmd, "end")
	}()

	fmt.Fprint(w, "running")
}

var deskCmd remoteCmd

func deskRun(actionName string) {

	action := getActionByName(actionName)

	logMsg := fmt.Sprintf("cmd: %s %s", action.Server, action.Cmd)

	deskCmd.lock.Lock()
	deskCmd.running = true
	go func() {
		defer deskCmd.lock.Unlock()
		println(logMsg, "start")
		ret, err := executeSSH(action.Cmd, action.Server)
		deskCmd.running = false
		if err == nil {
			voiceCmd.out = ret
			println(logMsg, " OK ")
		} else {
			println(logMsg, " FAIL ", err)
		}
	}()
}

//ViewDeskUp move desk up
func ViewDeskUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("refresh", "3;url=/")
	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRun("desk up")
		fmt.Fprint(w, "going up")
	}
}

//ViewDeskDown move desk down
func ViewDeskDown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("refresh", "3;url=/")

	if deskCmd.running {
		fmt.Fprint(w, "busy")
	} else {
		deskRun("desk down")
		fmt.Fprint(w, "going down")
	}
}
