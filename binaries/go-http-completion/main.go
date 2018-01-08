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

func completeRequests(partialName string, profileName string) {
	requestNames, err := profile.GetAvailableRequests(profileName)
	if err != nil {
		panic(err)
	}

	for _, requestName := range requestNames {
		if strings.HasPrefix(requestName, partialName) {
			fmt.Printf("@%s\n", requestName)
		}
	}
}

func main() {
	testString := os.Args[2]

	if strings.HasPrefix(testString, "+") {
		completeProfiles(testString[1:])
	} else if strings.HasPrefix(testString, "@") || strings.HasPrefix(testString, "\\@") {
		// Only complete request names if a profile is available
		if len(os.Args) == 4 && strings.HasPrefix(os.Args[3], "+") {
			partialName := testString[1:]
			if strings.HasPrefix(testString, "\\@") {
				partialName = testString[2:]
			}
			completeRequests(partialName, os.Args[3][1:])
		}
		fmt.Println("")
	}
}
