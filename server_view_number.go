package main

import (
	"fmt"
	"log"
	"net/http"
)

type viewNumberData struct {
	num int
}

var viewNumData viewNumberData

// ViewNumber increase number
func ViewNumber(w http.ResponseWriter, r *http.Request) {
	num := viewNumData.num

	log.Println("ViewNumber ", num)

	if num < 5 {
		viewNumData.num++
		w.Header().Set("refresh", "1")
	} else {
		w.Header().Set("refresh", "1;url=/")
	}

	fmt.Fprintf(w, "Hello %d ", num)
}
