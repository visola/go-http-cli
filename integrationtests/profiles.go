package integrationtests

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/visola/go-http-cli/profile"
	"github.com/visola/variables/variables"
)

// CreateProfile creates a profile file with the specified content in the current test
// profile directory
func CreateProfile(name string, content string) {
	profilesDir, err := profile.GetProfilesDir()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		os.Mkdir(profilesDir, 0777)
	}

	content = variables.ReplaceVariables(content, getContext())
	ioutil.WriteFile(path.Join(profilesDir, name+".yml"), []byte(content), 0777)
}
