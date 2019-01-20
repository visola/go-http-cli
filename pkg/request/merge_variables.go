package request

func mergeVariables(allVariables ...map[string]string) map[string]string {
	result := make(map[string]string)
	for i := len(allVariables) - 1; i >= 0; i-- {
		vars := allVariables[i]
		for key, val := range vars {
			result[key] = val
		}
	}
	return result
}
