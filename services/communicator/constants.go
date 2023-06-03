package communicator

const (
	// TODO: implement help command
	HelpCommand  = "!help"
	ExitCommand  = "!exit"
	ResetCommand = "!reset"
)

const cyclingChars = 9
const (
	initStatusText = "Chat initialization"
	generationText = "Generation"
)

type state int

const (
	startState state = iota
	initRunningState
	initCompletedState
	completionState
	completionDoneState
	errorState
)
