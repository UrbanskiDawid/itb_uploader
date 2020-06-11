package main
 
import (
    "fmt"
    "net/http"
)

var num int = 0;

func view_index(w http.ResponseWriter, r *http.Request) {
    num++;
    fmt.Fprintf(w, "Hello %d", num);
}

func main() {
    http.HandleFunc("/", view_index);    
    http.ListenAndServe(":8080",nil);
}