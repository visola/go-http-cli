package request

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/visola/go-http-cli/base"
	"github.com/visola/go-http-cli/profile"
)

// Request stores data required to configure a request to be executed
type Request struct {
	Body        string
	Cookies     []*http.Cookie
	Headers     map[string][]string
	Method      string
	QueryParams map[string][]string
	URL         string
}

// GetBody returns the body for this request
func (req *Request) GetBody() (string, error) {
	return req.Body, nil
}

// GetHeaders returns the headers for this request
func (req *Request) GetHeaders() map[string][]string {
	return req.Headers
}

// GetMethod returns the HTTP method for this request
func (req *Request) GetMethod() string {
	return req.Method
}

// LoadBodyFromFile loads data from a file and set it to the body, if not already set
func (req *Request) LoadBodyFromFile(fileName string) error {
	if fileName == "" {
		return nil
	}

	data, loadError := ioutil.ReadFile(fileName)
	if loadError != nil {
		return loadError
	}

	if req.Body != "" {
		return errors.New("Cannot set body and try to load from file at the same time")
	}

	req.Body = string(data)
	return nil
}

// Merge merges information from something compatible with a request into this request
func (req *Request) Merge(toMerge interface{}) error {
	withBody, ok := toMerge.(base.WithBody)
	if ok {
		body, err := withBody.GetBody()
		if err != nil {
			return err
		}
		req.MergeBody(body)
	}

	withHeader, ok := toMerge.(base.WithHeaders)
	if ok {
		req.MergeHeaders(withHeader.GetHeaders())
	}

	withMethod, ok := toMerge.(base.WithMethod)
	if ok {
		if withMethod.GetMethod() != "" {
			req.Method = withMethod.GetMethod()
		}
	}

	reqToMerge, ok := toMerge.(Request)
	if ok {
		req.Cookies = append(req.Cookies, reqToMerge.Cookies...)

		if reqToMerge.URL != "" {
			req.URL = reqToMerge.URL
		}
	}

	return nil
}

// MergeBody merges a body with this request
func (req *Request) MergeBody(toMerge string) {
	if toMerge != "" {
		req.Body = toMerge
	}
}

// MergeHeaders merges the passed headers into the request
func (req *Request) MergeHeaders(headers map[string][]string) {
	if req.Headers == nil {
		req.Headers = make(map[string][]string)
	}

	for header, values := range headers {
		req.Headers[header] = values
	}
}

// MergeOptions merges request options loaded from a profile
func (req *Request) MergeOptions(optionsToMerge profile.NamedRequest) {
	req.MergeBody(optionsToMerge.Body)
	req.MergeHeaders(optionsToMerge.Headers)

	if optionsToMerge.Method != "" {
		req.Method = optionsToMerge.Method
	}

	if optionsToMerge.URL != "" {
		req.URL = optionsToMerge.URL
	}
}
