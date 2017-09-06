package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testHandler)
	http.ListenAndServe(":8080", nil)
}