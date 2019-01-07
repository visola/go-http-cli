package request

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

// PostProcessSourceCode represents the source code to be executed
type PostProcessSourceCode struct {
	// SourceCode is the code to be executed
	SourceCode string

	// SourceFilePath is the path to the file where the code came from
	SourceFilePath string
}

// PostProcess processes the executed requests using the post processing script
func PostProcess(sourceCode PostProcessSourceCode, executedRequests []ExecutedRequestResponse, responseErr error) (string, error) {
	if sourceCode.SourceCode == "" {
		return "", nil
	}

	vm := otto.New()

	script, compileErr := vm.Compile(sourceCode.SourceFilePath, sourceCode.SourceCode)
	if compileErr != nil {
		return "", compileErr
	}

	output := ""
	printFunction := func(arg interface{}) {
		output += fmt.Sprint(arg)
	}

	vm.Set("executed", executedRequests)
	if len(executedRequests) > 0 {
		vm.Set("request", executedRequests[0].Request)
		vm.Set("response", executedRequests[0].Response)
		vm.Set("print", printFunction)
	}

	_, executeError := vm.Run(script)
	if executeError != nil {
		return "", executeError
	}

	return output, nil
}
