package variables

import "fmt"

// Variable is a named variable found in some string. It contains the name
// of the variable, associated tag if any and the position in the string where
// it was found
type Variable struct {
	End       int
	Name      string
	NameEnd   int
	NameStart int
	Start     int
	Tag       string
	TagEnd    int
	TagStart  int
}

// AsString returns a string representation of this variable
func (v Variable) AsString() string {
	if v.Tag == "" {
		return fmt.Sprintf("{%s}", v.Name)
	}
	return fmt.Sprintf("{%s:%s}", v.Name, v.Tag)
}
