package request

import "net/url"

func encodeValues(values map[string][]string) string {
	vals := url.Values{}
	for name, valuesForKey := range values {
		for _, value := range valuesForKey {
			vals.Add(name, value)
		}
	}
	return vals.Encode()
}
