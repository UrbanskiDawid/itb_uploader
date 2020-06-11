package main

import (
	"fmt"
	"net/http"
)

var num int = 0

var html string = `
<script>
</script>
`

func view_index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

func view_get(w http.ResponseWriter, r *http.Request) {
	num++

	if num < 10 {
		w.Header().Set("refresh", "1")
	} else {
		w.Header().Set("refresh", "1;url=/")
	}

	fmt.Fprintf(w, "Hello %d ", num)
}
