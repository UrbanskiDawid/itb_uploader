package views

import (
	"fmt"
	"net/http"
	"logging"
	"workers"
)

// ViewAllActions show all acions
func ViewAllActions(w http.ResponseWriter, r *http.Request) {
	logging.Log.Println("ViewAllActions")
	fmt.Println("Request ViewAllActions")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	actionNames := workers.GetActionNames()
	for _ ,name := range actionNames{
		fmt.Fprintf(w, "Action: %s</br>", name)
	} 

}

