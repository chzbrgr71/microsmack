package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Category string `json:"Category"`
	Item     string `json:"Item"`
	Value    string `json:"Value"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// gather values
	var hostname = getHostname()
	var appversion = "1.0.0"
	var gitSHA = os.Getenv("GIT_SHA")

	// call api for background color
	// https://blog.alexellis.io/golang-json-api-client
	// https://tutorialedge.net/golang/consuming-restful-api-with-go
	url := "http://localhost:8081/all"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(responseData))

	var configObject Config
	json.Unmarshal(responseData, &configObject)

	log.Printf("Category=%s", configObject.Category)
	log.Printf("Item=%s", configObject.Item)
	log.Printf("Value=%s", configObject.Value)

	// build web page
	var backColor = "green"
	var htmlHeader = "<!DOCTYPE html><html><font color=white><h1>Microsmack Homepage</h1>"
	fmt.Fprintf(w, htmlHeader)
	fmt.Fprintf(w, "<body style=background-color:%s>", backColor)
	fmt.Fprintf(w, "<p>Version: %s</p><p>Hostname: %s</p><p>Git: %s</p>", appversion, hostname, gitSHA)

}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "OK")
}

func getHostname() string {
	var result string
	localhostname, err := os.Hostname()

	if err != nil {
		result = "ERROR: Cannot find server hostname"
	} else {
		result = localhostname
	}
	return result
}
