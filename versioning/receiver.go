package versioning

// OldGreeter represents an old Receiver used by the server.
type OldGreeter struct{}

func (g OldGreeter) SayHello(str string) string {
	return "Hello " + str
}
