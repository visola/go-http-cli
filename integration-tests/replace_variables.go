package main

import "github.com/visola/variables/variables"

func replaceVariablesInArray(arrayIn ...string) []string {
	result := make([]string, len(arrayIn))
	for i, val := range arrayIn {
		result[i] = variables.ReplaceVariables(val, getContext())
	}
	return result
}
