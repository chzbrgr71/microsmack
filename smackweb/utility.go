package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getBackColor() string {
	// call api for background color
	url := "http://localhost:8081/all"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(responseData))

	var configObject Config
	json.Unmarshal(responseData, &configObject)
	//category := configObject.Category
	//item := configObject.Item
	value := configObject.Value

	return value

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
