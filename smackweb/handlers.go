package main

import (
	"fmt"
	"io"
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
	var gitSHA = os.Getenv("GIT_SHA")
	var appVersion = os.Getenv("APP_VERSION")
	var imageBuildDate = os.Getenv("IMAGE_BUILD_DATE")
	var kubeNodeName = os.Getenv("KUBE_NODE_NAME")
	var kubePodName = os.Getenv("KUBE_POD_NAME")
	var kubePodIP = os.Getenv("KUBE_POD_IP")

	// render page
	html := fmt.Sprintf("<!DOCTYPE html><html><font color=white><h1>Microsmack Cool Homepage</h1><body style=background-color:%s><p>Git: %s</p><p>App version: %s</p><p>Image build date: %s</p><p>Kubernetes node: %s</p><p>Kubernetes pod name: %s</p><p>Kubernetes pod IP: %s</p></body></html>", backColor, gitSHA, appVersion, imageBuildDate, kubeNodeName, kubePodName, kubePodIP)
	fmt.Fprintf(w, html)
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "RUNNING")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusBadGateway)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
