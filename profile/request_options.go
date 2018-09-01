package profile

import (
	"io/ioutil"
	"path/filepath"
)

// RequestOptions is a representation of a request that can be loaded from a profile.
type RequestOptions struct {
	Body         string
	FileToUpload string
	Headers      map[string][]string
	Method       string
	URL          string
	Values       map[string][]string
}

// GetBody returns the body for this RequestOptions
func (req RequestOptions) GetBody() (string, error) {
	if req.Body != "" {
		return req.Body, nil
	}

	// If there's a file to upload from profile, load it
	if req.FileToUpload != "" {
		path := req.FileToUpload

		// If not absolute, it's relative to the profiles dir
		if !filepath.IsAbs(path) {
			profileDir, profileDirError := GetProfilesDir()
			if profileDirError != nil {
				return "", profileDirError
			}
			path = filepath.Join(profileDir, req.FileToUpload)
		}

		data, loadErr := ioutil.ReadFile(path)
		if loadErr != nil {
			return "", loadErr
		}

		return string(data), nil
	}

	return "", nil
}

// GetHeaders returns the headers for this RequestOptions
func (req RequestOptions) GetHeaders() map[string][]string {
	return req.Headers
}

// GetValues returns the values for this RequestOptions
func (req RequestOptions) GetValues() map[string][]string {
	return req.Values
}
