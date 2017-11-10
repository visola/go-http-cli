package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/visola/go-http-cli/profile"
)

func completeProfiles(profilePrefix string) {
	profiles, err := profile.GetAvailableProfiles()
	if err != nil {
		panic(err)
	}

	for _, profile := range profiles {
		if strings.HasPrefix(profile, profilePrefix) {
			fmt.Printf("+%s\n", profile)
		}
	}
}

func main() {
	testString := os.Args[2]

	if strings.HasPrefix(testString, "+") {
		completeProfiles(testString[1:])
	}
}
