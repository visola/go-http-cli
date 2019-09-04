package session

import "github.com/visola/go-http-cli/pkg/model"

// SetVariableRequest is used to request a variable to be set
type SetVariableRequest struct {
	Values []model.KeyValuePair
}
