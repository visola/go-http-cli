package request

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

// PostProcess processes the executed requests using the post processing script
func PostProcess(options *ExecutionOptions, executedRequests []ExecutedRequestResponse, responseErr error) error {
	if options.PostProcessFile == "" {
		return nil
	}

	sourceCode, readError := ioutil.ReadFile(options.PostProcessFile)
	if readError != nil {
		return readError
	}

	vm := otto.New()

	script, compileErr := vm.Compile(options.PostProcessFile, sourceCode)
	if compileErr != nil {
		return compileErr
	}

	vm.Set("executed", executedRequests)
	if len(executedRequests) > 0 {
		vm.Set("request", executedRequests[0].Request)
		vm.Set("response", executedRequests[0].Response)
	}

	_, executeError := vm.Run(script)
	if executeError != nil {
		return executeError
	}

	return nil
}
