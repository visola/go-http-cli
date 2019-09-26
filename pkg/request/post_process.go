package request

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/op/go-logging"
	"github.com/robertkrimen/otto"
	"github.com/visola/go-http-cli/pkg/profile"
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
	sourceCode := executedRequests[len(executedRequests)-1].Request.PostProcessCode
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

func createAddRequestFunction(context *PostProcessContext, executionContext *ExecutionContext) func(otto.Value) {
	return func(value otto.Value) {
		unconfiguredRequest := Request{}
		requestName := ""

		if value.IsString() {
			valueAsString := otto.Value.String(value)
			if strings.HasPrefix(valueAsString, "@") {
				requestName = valueAsString[1:]
			} else {
				unconfiguredRequest.URL = valueAsString
			}
		}

		if value.IsObject() {
			toAddAsMap, err := value.Export()
			if err != nil {
				context.Output += fmt.Sprintf("Error while converting object to map: %s\n%s\n", value, err.Error())
				log.Error("Error while converting request object to map.", err)
				return
			}

			var toAdd Request
			if err := mapstructure.Decode(toAddAsMap, &toAdd); err != nil {
				context.Output += fmt.Sprintf("Error while converting map to request object.\n%s\n", err.Error())
				log.Error("Error while converting map to request object.", err)
				return
			}

			unconfiguredRequest = toAdd
		}

		// Profiles were loaded before, no need to handle error here
		mergedProfile, _ := profile.LoadAndMergeProfiles(executionContext.ProfileNames)

		configuredRequest, err := ConfigureRequestSimple(unconfiguredRequest, &mergedProfile, requestName)
		if err != nil {
			context.Output += fmt.Sprintf("Error while configuring request: %s\n", err.Error())
			return
		}

		context.Requests = append(context.Requests, *configuredRequest)
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
	vm.Set("addRequest", createAddRequestFunction(context, executionContext))
	vm.Set("print", createPrintFunction(context))
	vm.Set("println", createPrintlnFunction(context))

	vm.Set("executed", executedRequests)
	if len(executedRequests) > 0 {
		vm.Set("request", executedRequests[0].Request)
		vm.Set("response", executedRequests[0].Response)
	}

	return context
}
