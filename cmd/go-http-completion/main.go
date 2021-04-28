package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/visola/go-http-cli/pkg/profile"
)

func bashCompletion() {
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

func getActiveProfiles() []string {
	profiles, err := profile.GetAvailableProfiles()
	if err != nil {
		panic(err)
	}

	profilesMap := make(map[string]bool)
	for _, p := range profiles {
		profilesMap[p] = true
	}

	active := make([]string, 0)
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "+") {
			p := arg[1:]
			if _, exists := profilesMap[p]; exists {
				active = append(active, p)
			}
		}
	}
	return active
}

func listProfiles() {
	profiles, err := profile.GetAvailableProfiles()
	if err != nil {
		panic(err)
	}

	for i, profile := range profiles {
		profiles[i] = "+" + profile
	}
	fmt.Print(strings.Join(profiles, " "))
}

func listRequests(activeProfiles []string) {
	requests := make([]string, 0)
	for _, activeProfile := range activeProfiles {
		requestNames, err := profile.GetAvailableRequests(activeProfile)
		if err != nil {
			panic(err)
		}

		requests = append(requests, requestNames...)
	}

	for i, req := range requests {
		requests[i] = "@" + req
	}
	fmt.Print(strings.Join(requests, " "))
}

func shouldListRequests() bool {
	for _, arg := range os.Args {
		if arg == "--requests" {
			return true
		}
	}
	return false
}

func zshCompletion() {
	fmt.Print(`#compdef _http http

function _http () { \
	profiles=$(go-http-completion --profiles)
	requests=$(go-http-completion --requests $words)
	
	_arguments -s \
		'-X+[HTTP method to use]:method:(GET POST PUT DELETE)' \
		'-T+[Use file as body]:upload file:_files' \
		"::profile:($profiles)" \
		"::method:($requests)"
}`)
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--profiles" {
		listProfiles()
		return
	}

	if len(os.Args) == 2 && os.Args[1] == "zsh" {
		zshCompletion()
		return
	}

	if shouldListRequests() {
		listRequests(getActiveProfiles())
		return
	}

	if len(os.Args) == 1 {
		panic("Should not run this without arguments")
	}

	bashCompletion()
}
