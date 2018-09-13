package profile

// GetAvailableRequests returns a list of all the names of requests available in a profile.
func GetAvailableRequests(profileName string) ([]string, error) {
	options, err := LoadProfile(profileName)

	if err != nil {
		return nil, err
	}

	requestNames := make([]string, 0)

	for requestName := range options.NamedRequest {
		requestNames = append(requestNames, requestName)
	}

	return requestNames, nil
}
