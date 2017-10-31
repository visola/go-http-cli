package request

import (
	"fmt"
)

func ExampleParseURL() {
	base := "http://localhost:3000/api/v1"
	url := "/${companyId}/contacts"
	context := map[string]string{
		"companyId": "123456",
	}
	fmt.Println(ParseURL(base, url, context))
	// Output: http://localhost:3000/api/v1/123456/contacts
}
