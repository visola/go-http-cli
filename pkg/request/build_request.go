package request

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/visola/go-http-cli/pkg/ioutil"
)

// BuildRequest builds an http.Request from a configured request.Request
func BuildRequest(processedRequest Request) (*http.Request, error) {
	parsedURL, urlError := url.Parse(processedRequest.URL)
	if urlError != nil {
		return nil, urlError
	}

	rawQueryPieces := make([]string, 0)
	if parsedURL.RawQuery != "" {
		rawQueryPieces = append(rawQueryPieces, parsedURL.RawQuery)
	}

	encodedQueryFromValues := encodeValues(processedRequest.QueryParams)
	if encodedQueryFromValues != "" {
		rawQueryPieces = append(rawQueryPieces, encodedQueryFromValues)
	}

	parsedURL.RawQuery = strings.Join(rawQueryPieces, "&")

	req, reqErr := http.NewRequest(processedRequest.Method, parsedURL.String(), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	for _, cookie := range processedRequest.Cookies {
		req.AddCookie(cookie)
	}

	for k, vs := range processedRequest.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	req.Body = ioutil.CreateCloseableBufferString(processedRequest.Body)
	return req, nil
}
