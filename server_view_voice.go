package main

import (
	"log"
	"net/http"
	"sync"
)

type viewVoiceMemory struct {
	running bool
	lock    sync.Mutex
	out     string
}

var voiceCmd viewVoiceMemory

//ViewVoice play voice via SSH
func ViewVoice(w http.ResponseWriter, r *http.Request) {
	log.Println("ViewVoice")
	actionName := "voice"

	voiceCmd.lock.Lock()
	voiceCmd.running = true
	go func() {
		defer voiceCmd.lock.Unlock()
		println("voice begin")
		ret, retErr, err := executeAction(actionName)
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
