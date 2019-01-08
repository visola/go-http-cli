package request

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

// PostProcessContext stores
type PostProcessContext struct {
	Output string
}

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

	context := preparePostProcessContext(vm, executedRequests, responseErr)

	_, executeError := vm.Run(script)
	if executeError != nil {
		return "", executeError
	}

	return context.Output, nil
}

func createPrintFunction(context *PostProcessContext) func(...interface{}) {
	return func(args ...interface{}) {
		context.Output += fmt.Sprint(args...)
	}
}

func createPrintlnFunction(context *PostProcessContext) func(...interface{}) {
	return func(args ...interface{}) {
		context.Output += fmt.Sprint(args...) + "\n"
	}
}

func preparePostProcessContext(vm *otto.Otto, executedRequests []ExecutedRequestResponse, responseErr error) *PostProcessContext {
	context := &PostProcessContext{
		Output: "",
	}

	vm.Set("print", createPrintFunction(context))
	vm.Set("println", createPrintlnFunction(context))

	vm.Set("executed", executedRequests)
	if len(executedRequests) > 0 {
		vm.Set("request", executedRequests[0].Request)
		vm.Set("response", executedRequests[0].Response)
	}

	return context
}
