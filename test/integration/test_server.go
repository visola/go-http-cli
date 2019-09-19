package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Request stores a request that was received by the test server
type Request struct {
	Body    string
	Cookies []*http.Cookie
	Headers map[string][]string
	Method  string
	Path    string
	Query   map[string][]string
}

// ReplyWith gives specifications the ability to ask the server to reply in a specific way
type ReplyWith struct {
	Headers map[string][]string
	Body string
}

const defaultBody = "Hello world!"

var allRequests = make([]Request, 0)
var lastRequest Request
var testServer *httptest.Server

var replyWith = defaultReplyWith()

func startTestServer() {
	allRequests = make([]Request, 0)
	testServer = httptest.NewServer(http.HandlerFunc(handleRequest))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}

	lastRequest = Request{
		Body:    string(body),
		Cookies: r.Cookies(),
		Headers: toLowerCaseHeaders(r.Header),
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.Query(),
	}

	allRequests = append(allRequests, lastRequest)

	for header, values := range replyWith.Headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// TODO - Store request received
	fmt.Fprintln(w, replyWith.Body);

	// Clean up reply with after finished
	replyWith = defaultReplyWith()
}

func defaultReplyWith() ReplyWith {
	return ReplyWith{
		Body: defaultBody,
	}
}

func prepareReply(r ReplyWith) {
	replyWith = r
}

func toLowerCaseHeaders(header http.Header) map[string][]string {
	result := make(map[string][]string)
	for header, values := range header {
		result[strings.ToLower(header)] = values
	}
	return result
}
