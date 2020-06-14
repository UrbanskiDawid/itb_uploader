package views

import (
	"net/http"
	"sync"
	"fmt"
	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/UrbanskiDawid/itb_uploader/actions"
)

type viewVoiceMemory struct {
	running bool
	lock    sync.Mutex
	out     string
}

var voiceCmd viewVoiceMemory

//ViewVoice play voice via SSH
func ViewVoice(w http.ResponseWriter, r *http.Request) {
	logging.Log.Println("ViewVoice")
	fmt.Println("Request ViewVoice")

	actionName := "voice"

	voiceCmd.lock.Lock()
	voiceCmd.running = true
	go func() {
		defer voiceCmd.lock.Unlock()
		println("voice begin")
		ret, retErr, err := actions.ExecuteAction(actionName)
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
