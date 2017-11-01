// Package request contains all the things used to create requests.
package request

import (
	"bytes"
	"net/http"

	"github.com/visola/go-http-cli/config"
)

type bodyBuffer struct {
	*bytes.Buffer
}

func (bb *bodyBuffer) Close() error {
	return nil
}

// BuildRequest builds a Request from a Configuration.
func BuildRequest(configuration config.Configuration) (*http.Request, error) {
	url := ParseURL(configuration.BaseURL(), configuration.URL(), configuration.Variables())

	req, reqErr := http.NewRequest(configuration.Method(), url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	for k, vs := range configuration.Headers() {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	if configuration.Body() != "" {
		req.Body = &bodyBuffer{bytes.NewBufferString(configuration.Body())}
	}

	return req, nil
}
