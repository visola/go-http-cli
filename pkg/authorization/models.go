package authorization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Authorization represents an HTTP authorization
type Authorization struct {
	AuthorizationType string
	Password          string
	Token             string
	Username          string
}

// IsValid checks if this authorization is valid or not
func (auth Authorization) IsValid() error {
	authType := strings.ToLower(auth.AuthorizationType)

	if authType == "basic" {
		if auth.Username == "" || auth.Password == "" {
			return fmt.Errorf("Username and password must not be empty but where '%s' and '%s' respectively", auth.Username, auth.Password)
		}

		return nil
	}

	if authType == "bearer" {
		if auth.Token == "" {
			return errors.New("Token must not be empty for Bearer auth")
		}

		return nil
	}

	return fmt.Errorf("Unsupported auth type: %s", authType)
}

// ToHeaderKey returns the key to be used when adding this authorization as a header.
func (auth Authorization) ToHeaderKey() string {
	return "Authorization"
}

// ToHeaderValue returns the value to be used when adding this authorization as a header.
func (auth Authorization) ToHeaderValue() (string, error) {
	validationError := auth.IsValid()

	if validationError != nil {
		return "", validationError
	}

	authType := strings.ToLower(auth.AuthorizationType)

	if authType == "basic" {
		toEncode := fmt.Sprintf("%s:%s", auth.Username, auth.Password)
		return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(toEncode))), nil
	}

	if authType == "bearer" {
		return fmt.Sprintf("Bearer %s", auth.Token), nil
	}

	return "", fmt.Errorf("Unsupported auth type: %s", authType)
}
