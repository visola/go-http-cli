package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthorization(t *testing.T) {
	authorization := Authorization{
		AuthorizationType: "Basic",
		Password:          "myPassword",
		Username:          "myUsername",
	}

	encoded, encodeError := authorization.ToHeaderValue()

	assert.Nil(t, encodeError, "Should encode correctly using Basic auth")
	assert.Equal(t, "Basic bXlVc2VybmFtZTpteVBhc3N3b3Jk", encoded, "Should generate correct header value")
}

func TestBearerAuthorization(t *testing.T) {
	token := "my-very-long-token"

	authorization := Authorization{
		AuthorizationType: "Bearer",
		Token:             token,
	}

	encoded, encodeError := authorization.ToHeaderValue()

	assert.Nil(t, encodeError, "Should encode correctly using Basic auth")
	assert.Equal(t, "Bearer "+token, encoded, "Should generate correct header value")
}
