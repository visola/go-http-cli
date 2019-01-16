package profile

import "fmt"

// FindNamedRequest finds a named request in a group of already merged profiles
func FindNamedRequest(mergedProfile *Options, requestName string) (NamedRequest, error) {
	if requestName == "" {
		return NamedRequest{}, nil
	}

	var namedRequest NamedRequest

	var exists bool
	if namedRequest, exists = mergedProfile.NamedRequest[requestName]; requestName != "" && !exists {
		return NamedRequest{}, fmt.Errorf("Request with name %s not found", requestName)
	}

	return namedRequest, nil
}
