package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// Request stores a request that was received by the test server
type Request struct {
	Body    string
	Headers map[string][]string
	Method  string
	Path    string
}

var lastRequest Request
var testServer *httptest.Server

func startTestServer() {
	testServer = httptest.NewServer(http.HandlerFunc(handleRequest))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}

	lastRequest = Request{
		Body:    string(body),
		Headers: r.Header,
		Method:  r.Method,
		Path:    r.URL.Path,
	}

	// TODO - Store request received
	// TODO - Return something useful
	fmt.Fprintln(w, "Hello world!")
}
