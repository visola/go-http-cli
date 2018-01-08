package profile

import (
	"fmt"
)

// LoadRequestOptions loads request configurations from the profile names passed in. It returns
// an error if it doesn't find a request configuration with the specified name.
func LoadRequestOptions(requestNameToFind string, profiles []string) (*RequestOptions, error) {
	for _, profileName := range profiles {
		profileOptions, err := LoadProfile(profileName)
		if err != nil {
			return nil, err
		}

		var result *RequestOptions
		for requestName, requestConfiguration := range profileOptions.RequestOptions {
			if requestNameToFind == requestName {
				result = &requestConfiguration
			}
		}

		// Return the last one found
		if result != nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("Request '%s' not found in profiles '%s'", requestNameToFind, profiles)
}
