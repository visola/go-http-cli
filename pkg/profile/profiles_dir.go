package profile

import (
	"os"
	"os/user"
)

const profilesDirEnvVariable = "GO_HTTP_PROFILES"

// GetProfilesDir return the directory where profiles are stored
func GetProfilesDir() (string, error) {
	profilesDir := os.Getenv(profilesDirEnvVariable)
	if profilesDir == "" {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		profilesDir = user.HomeDir + "/go-http-cli"
	}
	return profilesDir, nil
}
