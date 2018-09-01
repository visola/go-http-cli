package request

import (
	"net/http"

	"github.com/visola/go-http-cli/base"
	"github.com/visola/go-http-cli/profile"
)

// Request stores data required to configure a request to be executed
type Request struct {
	Body    string
	Cookies []*http.Cookie
	Headers map[string][]string
	Method  string
	URL     string
	Values  map[string][]string
}

// GetBody returns the body for this request
func (req *Request) GetBody() (string, error) {
	return req.Body, nil
}

// GetHeaders returns the headers for this request
func (req *Request) GetHeaders() map[string][]string {
	return req.Headers
}

// GetValues returns the values for this request
func (req *Request) GetValues() map[string][]string {
	return req.Values
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

	withValues, ok := toMerge.(base.WithValues)
	if ok {
		req.MergeValues(withValues.GetValues())
	}

	reqToMerge, ok := toMerge.(Request)
	if ok {
		req.Cookies = append(req.Cookies, reqToMerge.Cookies...)

		if reqToMerge.Method != "" {
			req.Method = reqToMerge.Method
		}

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
func (req *Request) MergeOptions(optionsToMerge profile.RequestOptions) {
	req.MergeBody(optionsToMerge.Body)
	req.MergeHeaders(optionsToMerge.Headers)
	req.MergeValues(optionsToMerge.Values)

	if optionsToMerge.Method != "" {
		req.Method = optionsToMerge.Method
	}

	if optionsToMerge.URL != "" {
		req.URL = optionsToMerge.URL
	}
}

// MergeValues merges values into this request
func (req *Request) MergeValues(valuesToMerge map[string][]string) {
	if req.Values == nil {
		req.Values = make(map[string][]string)
	}

	for name, values := range valuesToMerge {
		req.Values[name] = values
	}
}
