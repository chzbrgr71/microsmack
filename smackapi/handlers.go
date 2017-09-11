package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Key      string `json:"Key"`
	Category string `json:"Category"`
	Item     string `json:"Item"`
	Value    string `json:"Value"`
}

type Configs []Config

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RUNNING")
}

func returnAllConfigs(w http.ResponseWriter, r *http.Request) {
	configs := Configs{
		Config{Key: "1", Category: "UI", Item: "Background Color", Value: "Blue"},
		Config{Key: "2", Category: "k8s", Item: "Kubernetes Node", Value: "minkube"},
		Config{Key: "3", Category: "k8s", Item: "Kubernetes Pod", Value: "pod-name"},
		Config{Key: "4", Category: "k8s", Item: "Kubernetes IP", Value: "192.168.1.1"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(configs); err != nil {
		panic(err)
	}
}

func returnColor(w http.ResponseWriter, r *http.Request) {
	configs := Config{Key: "1", Category: "UI", Item: "Background Color", Value: "Blue"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(configs); err != nil {
		panic(err)
	}
}

func returnSingleConfig(w http.ResponseWriter, r *http.Request) {
	// need to implement this... return single item based on key
	vars := mux.Vars(r)
	key := vars["key"]

	fmt.Fprintf(w, "Key: "+key)
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "RUNNING")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check. can simulate error with http status
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusBadGateway)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
