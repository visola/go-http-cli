package request

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/visola/go-http-cli/pkg/base"
	"github.com/visola/go-http-cli/pkg/profile"
)

// Request stores data required to configure a request to be executed
type Request struct {
	AllowInsecure   bool
	Body            string
	Cookies         []*http.Cookie
	Headers         map[string][]string
	Method          string
	PostProcessCode PostProcessSourceCode
	QueryParams     map[string][]string
	URL             string
}

// GetAllowInsecure returns if this named request allow insecure HTTP connections
func (req Request) GetAllowInsecure() bool {
	return req.AllowInsecure
}

// GetBody returns the body for this request
func (req Request) GetBody() (string, error) {
	return req.Body, nil
}

// GetHeaders returns the headers for this request
func (req Request) GetHeaders() map[string][]string {
	return req.Headers
}

// GetMethod returns the HTTP method for this request
func (req Request) GetMethod() string {
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
	if withBody, ok := toMerge.(base.WithBody); ok {
		body, err := withBody.GetBody()
		if err != nil {
			return err
		}
		req.MergeBody(body)
	}

	if withAllowInsecure, ok := toMerge.(base.WithAllowInsecure); ok {
		req.AllowInsecure = req.AllowInsecure || withAllowInsecure.GetAllowInsecure()
	}

	if withHeader, ok := toMerge.(base.WithHeaders); ok {
		req.MergeHeaders(withHeader.GetHeaders())
	}

	if withMethod, ok := toMerge.(base.WithMethod); ok {
		if withMethod.GetMethod() != "" {
			req.Method = withMethod.GetMethod()
		}
	}

	if reqToMerge, ok := toMerge.(Request); ok {
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
