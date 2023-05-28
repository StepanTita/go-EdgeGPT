package communicator

const (
	HelpCommand  = "!help"
	ExitCommand  = "!exit"
	ResetCommand = "!reset"
)

const cyclingChars = 9
const (
	initStatusText = "Chat initialization"
	generationText = "Generation"
)

const markdownPrefix = "Format the response as Markdown."

type state int

const (
	startState state = iota
	initRunningState
	initCompletedState
	completionState
	completionDoneState
	errorState
)
