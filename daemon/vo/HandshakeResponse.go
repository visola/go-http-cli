package vo

// HandshakeResponse is the response sent by the daemon when someone is checking if it's up.
type HandshakeResponse struct {
	MajorVersion int8
	MinorVersion int8
}
