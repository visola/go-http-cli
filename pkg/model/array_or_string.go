package model

// ArrayOrString is a type used to unmarshal a field as an array of strings or as a single string from YAMLs
type ArrayOrString []string

// UnmarshalYAML implement the unmarshal from YAML package
func (v *ArrayOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		// Apparently value is not an array
		var single string
		err := unmarshal(&single)
		if err != nil {
			// Still can't parse it
			return err
		}
		*v = []string{single}
	} else {
		*v = multi
	}
	return nil
}

// ToMapOfArrayOfStrings transforms a map with string keys into a map of array of strings
func ToMapOfArrayOfStrings(original map[string]ArrayOrString) map[string][]string {
	result := make(map[string][]string)
	for key, values := range original {
		result[key] = values
	}
	return result
}
