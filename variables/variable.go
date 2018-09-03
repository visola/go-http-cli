package variables

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
