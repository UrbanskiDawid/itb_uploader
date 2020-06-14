package views

import (
	"fmt"
	"net/http"
	logging "github.com/UrbanskiDawid/itb_uploader/logging"
)

type viewNumberData struct {
	num int
}

var viewNumData viewNumberData

// ViewNumber increase number
func ViewNumber(w http.ResponseWriter, r *http.Request) {
	num := viewNumData.num

	logging.Log.Println("ViewNumber ", num)
	fmt.Println("Request ViewNumber")
	

	if num < 5 {
		viewNumData.num++
		w.Header().Set("refresh", "1")
	} else {
		w.Header().Set("refresh", "1;url=/")
	}

	fmt.Fprintf(w, "Hello %d ", num)
}
