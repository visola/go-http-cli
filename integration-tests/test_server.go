package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

var testServer *httptest.Server

func startTestServer() {
	testServer = httptest.NewServer(http.HandlerFunc(handleRequest))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// TODO - Store request received
	// TODO - Return something useful
	fmt.Fprintln(w, "Hello world!")
}
