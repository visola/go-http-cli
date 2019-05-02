package integration

func getContext() map[string]string {
	return map[string]string{
		"test-server": testServer.URL,
	}
}
