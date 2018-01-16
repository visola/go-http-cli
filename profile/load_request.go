package profile

import (
	"fmt"
)

// LoadRequestOptions loads request configurations from the profile names passed in. It returns
// an error if it doesn't find a request configuration with the specified name. It will return the
// last request options found, if the same name is found in multiple profiles.
func LoadRequestOptions(requestNameToFind string, profiles []string) (*RequestOptions, error) {
	result := make([]RequestOptions, 0)
	for _, profileName := range profiles {
		profileOptions, err := LoadProfile(profileName)
		if err != nil {
			return nil, err
		}

		for requestName, requestConfiguration := range profileOptions.RequestOptions {
			if requestNameToFind == requestName {
				result = append(result, requestConfiguration)
			}
		}
	}

	if len(result) > 0 {
		return &result[len(result)-1], nil
	}

	return nil, fmt.Errorf("Request '%s' not found in profiles '%s'", requestNameToFind, profiles)
}
