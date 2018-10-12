package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Request stores a request that was received by the test server
type Request struct {
	Method string
}

var lastRequest Request
var testServer *httptest.Server

func startTestServer() {
	testServer = httptest.NewServer(http.HandlerFunc(handleRequest))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	lastRequest = Request{
		Method: r.Method,
	}

	// TODO - Store request received
	// TODO - Return something useful
	fmt.Fprintln(w, "Hello world!")
}
