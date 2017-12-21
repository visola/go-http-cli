package profile

// We want to accept headers as single string or array of strings
type arrayOrString []string

func (v *arrayOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
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

func toMapOfArrayOfStrings(original map[string]arrayOrString) map[string][]string {
	result := make(map[string][]string)
	for key, values := range original {
		result[key] = values
	}
	return result
}
