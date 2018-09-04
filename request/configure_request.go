package request

import (
	"fmt"

	"github.com/visola/go-http-cli/profile"
)

// ConfigureRequest configures a request to be executed based on the provided options
func ConfigureRequest(unconfiguredRequest Request, requestName string, profileNames []string) (*Request, error) {
	mergedProfile, profileError := loadAndMergeProfiles(profileNames)
	if profileError != nil {
		return nil, profileError
	}

	requestOptions, requestOptionsErr := findNamedRequest(mergedProfile, requestName)
	if requestOptionsErr != nil {
		return nil, requestOptionsErr
	}

	configuredRequest := Request{}

	configuredRequest.Merge(mergedProfile)
	configuredRequest.Merge(requestOptions)
	configuredRequest.Merge(unconfiguredRequest)

	configuredRequest.URL = ParseURL(mergedProfile.BaseURL, configuredRequest.URL, requestOptions.URL)

	return &configuredRequest, nil
}

func findNamedRequest(mergedProfile profile.Options, requestName string) (profile.RequestOptions, error) {
	if requestName == "" {
		return profile.RequestOptions{}, nil
	}

	var requestOptions profile.RequestOptions

	var exists bool
	if requestOptions, exists = mergedProfile.RequestOptions[requestName]; requestName != "" && !exists {
		return profile.RequestOptions{}, fmt.Errorf("Request with name %s not found", requestName)
	}

	return requestOptions, nil
}

func loadAndMergeProfiles(profileNames []string) (profile.Options, error) {
	profiles := make([]profile.Options, len(profileNames))

	for index, profileName := range profileNames {
		loadedProfile, profileError := profile.LoadProfile(profileName)
		if profileError != nil {
			return profile.Options{}, profileError
		}
		profiles[index] = loadedProfile
	}

	return profile.MergeOptions(profiles), nil
}
