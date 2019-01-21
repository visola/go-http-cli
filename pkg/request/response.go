package request

// Response is the response from the daemon after executing a request
type Response struct {
	Body       string
	Headers    map[string][]string
	Protocol   string
	StatusCode int
	Status     string
}
