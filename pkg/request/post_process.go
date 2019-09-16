package request

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/op/go-logging"
	"github.com/robertkrimen/otto"
	"github.com/visola/go-http-cli/pkg/session"
)

var log = logging.MustGetLogger("post-processor")

// PostProcessContext stores
type PostProcessContext struct {
	Output   string
	Requests []Request
}

// PostProcessSourceCode represents the source code to be executed
type PostProcessSourceCode struct {
	// SourceCode is the code to be executed
	SourceCode string

	// SourceFilePath is the path to the file where the code came from
	SourceFilePath string
}

// PostProcess processes the executed requests using the post processing script
func PostProcess(executionContext *ExecutionContext, executedRequests []ExecutedRequestResponse, responseErr error) (*PostProcessContext, error) {
	sourceCode := executionContext.PostProcessCode
	if sourceCode.SourceCode == "" {
		return &PostProcessContext{}, nil
	}

	vm := otto.New()

	script, compileErr := vm.Compile(sourceCode.SourceFilePath, sourceCode.SourceCode)
	if compileErr != nil {
		return &PostProcessContext{}, compileErr
	}

	context := preparePostProcessContext(vm, executionContext, executedRequests, responseErr)

	_, executeError := vm.Run(script)
	if executeError != nil {
		return &PostProcessContext{}, executeError
	}

	return context, nil
}

func createAddVariableFunction(executionContext *ExecutionContext) func(string, string) {
	return func(name string, value string) {
		session.SetVariable(executionContext.Session.Host, name, value)
	}
}

func createAddRequestFunction(context *PostProcessContext) func(otto.Value) {
	return func(value otto.Value) {
		if value.IsString() {
			context.Requests = append(context.Requests, Request{
				URL: otto.Value.String(value),
			})
			return
		}

		if value.IsObject() {
			toAddAsMap, err := value.Export()
			if err != nil {
				log.Error("Error while converting request object to map.", err)
				return
			}

			var toAdd Request
			if err := mapstructure.Decode(toAddAsMap, &toAdd); err != nil {
				log.Error("Error while converting map to request object.", err)
				return
			}

			context.Requests = append(context.Requests, toAdd)
		}
	}
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

func preparePostProcessContext(vm *otto.Otto, executionContext *ExecutionContext, executedRequests []ExecutedRequestResponse, responseErr error) *PostProcessContext {
	context := &PostProcessContext{
		Output:   "",
		Requests: make([]Request, 0),
	}

	vm.Set("addVariable", createAddVariableFunction(executionContext))
	vm.Set("addRequest", createAddRequestFunction(context))
	vm.Set("print", createPrintFunction(context))
	vm.Set("println", createPrintlnFunction(context))

	vm.Set("executed", executedRequests)
	if len(executedRequests) > 0 {
		vm.Set("request", executedRequests[0].Request)
		vm.Set("response", executedRequests[0].Response)
	}

	return context
}
