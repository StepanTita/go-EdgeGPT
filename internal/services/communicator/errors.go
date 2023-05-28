package communicator

type communicatorError struct {
	err    error
	reason string
}

func (m communicatorError) Error() string {
	return m.err.Error()
}
