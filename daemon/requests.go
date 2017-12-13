package daemon

import (
	"github.com/visola/go-http-cli/options"
)

// ExecuteRequestRequest request to be sent to daemon when a request needs to be executed.
type ExecuteRequestRequest struct {
	Options *options.CommandLineOptions
}
