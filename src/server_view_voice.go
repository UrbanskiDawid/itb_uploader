package main

import (
	"net/http"
	"sync"
	workers "dawidurbanski.pl/itb_uploader/workers"
)

type viewVoiceMemory struct {
	running bool
	lock    sync.Mutex
	out     string
}

var voiceCmd viewVoiceMemory

//ViewVoice play voice via SSH
func ViewVoice(w http.ResponseWriter, r *http.Request) {
	Log.Println("ViewVoice")
	actionName := "voice"

	voiceCmd.lock.Lock()
	voiceCmd.running = true
	go func() {
		defer voiceCmd.lock.Unlock()
		println("voice begin")
		ret, retErr, err := workers.ExecuteAction(actionName)
		voiceCmd.running = false
		if err == nil {
			voiceCmd.out = ret
			println("voice end: OK", ret, retErr)
		} else {
			println("voice end: FAIL", retErr, err)
		}
	}()

	w.Header().Set("refresh", "1;url=/")
}
