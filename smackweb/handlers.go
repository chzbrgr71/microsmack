package main

import (
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Category string `json:"Category"`
	Item     string `json:"Item"`
	Value    string `json:"Value"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// gather values
	var backColor = getBackColor()
	var hostname = getHostname()
	var appversion = "1.0.0"
	var gitSHA = os.Getenv("GIT_SHA")

	// render page
	html := fmt.Sprintf("<!DOCTYPE html><html><font color=white><h1>Microsmack Homepage</h1><body style=background-color:%s><p>Version: %s</p><p>Hostname: %s</p><p>Git: %s</p></body></html>", backColor, appversion, hostname, gitSHA)
	fmt.Fprintf(w, html)
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "RUNNING")
}
