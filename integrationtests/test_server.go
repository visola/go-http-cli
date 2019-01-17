package integrationtests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Request stores a request that was received by the test server
type Request struct {
	Body    string
	Headers map[string][]string
	Method  string
	Path    string
	Query   map[string][]string
}

// ReplyWith gives specifications the ability to ask the server to reply in a specific way
type ReplyWith struct {
	Headers map[string][]string
}

// WrapWithTestServer initializes the test server and make sure it will tear down correctly after
func WrapWithTestServer(toWrap func(*testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		startTestServer()
		defer testServer.Close()

		toWrap(t)
	}
}

var lastRequest Request
var replyWith ReplyWith
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
		Query:   r.URL.Query(),
	}

	for header, values := range replyWith.Headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// TODO - Store request received
	// TODO - Return something useful
	fmt.Fprintln(w, "Hello world!")

	// Clean up reply with after finished
	replyWith = ReplyWith{}
}

func prepareReply(r ReplyWith) {
	replyWith = r
}
