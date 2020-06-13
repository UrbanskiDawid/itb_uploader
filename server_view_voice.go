package main

import (
	"fmt"
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
