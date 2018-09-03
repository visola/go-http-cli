package variables

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindVariables(t *testing.T) {
	t.Run("Find variable without tag", testFindVariableWithoutTag)
	t.Run("Find variable with tag", testFindVariableWithTag)
	t.Run("Find many mixed variables", testFindManyMixedVariables)
}

func testFindVariableWithoutTag(t *testing.T) {
	stringWithVariable := "The man called {name} was found living in a box."
	vars := FindVariables(stringWithVariable)

	assert.Equal(t, 1, len(vars), "Should find one variable")
	if len(vars) != 1 {
		return
	}
	variable := vars[0]
	assert.Equal(t, "name", variable.Name, "Should find variable with correct name")
	assert.Equal(t, "", variable.Tag, "Should not find any tag")
	assert.Equal(t, 15, variable.Start, "Should find start position")
	assert.Equal(t, 21, variable.End, "Should find end position")
	assert.Equal(t, 16, variable.NameStart, "Should find name start position")
	assert.Equal(t, 20, variable.NameEnd, "Should find name end position")
}

func testFindVariableWithTag(t *testing.T) {
	stringWithVariable := "The man called {name:required} was found living in a box."
	vars := FindVariables(stringWithVariable)

	assert.Equal(t, 1, len(vars), "Should find one variable")
	if len(vars) != 1 {
		return
	}

	variable := vars[0]
	assert.Equal(t, "name", vars[0].Name, "Should find variable with correct name")
	assert.Equal(t, "required", vars[0].Tag, "Should find correct tag")
	assert.Equal(t, 15, variable.Start, "Should find start position")
	assert.Equal(t, 30, variable.End, "Should find end position")
	assert.Equal(t, 16, variable.NameStart, "Should find name start position")
	assert.Equal(t, 20, variable.NameEnd, "Should find name end position")
	assert.Equal(t, 21, variable.TagStart, "Should find tag start position")
	assert.Equal(t, 29, variable.TagEnd, "Should find tag end position")
}

func testFindManyMixedVariables(t *testing.T) {
	jsonWithVariables := `
		{
			"name": "{name}",
			"password": "{password:request}",
			"company": "{company}",
			"employeeId": {employeeId:number}
		}
	`

	vars := FindVariables(jsonWithVariables)

	assert.Equal(t, 4, len(vars), "Should find all variables")
}
