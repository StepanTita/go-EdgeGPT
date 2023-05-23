package communicator

func checkExit(in string, breakline bool) bool {
	return in == ExitCommand && breakline
}
