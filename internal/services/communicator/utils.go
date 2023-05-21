package communicator

import "github.com/c-bata/go-prompt"

func executor(t string) {
	// TODO: implement here conversation logic
	return
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: HelpCommand},
		{Text: ExitCommand},
		{Text: ResetCommand},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func checkExit(in string, breakline bool) bool {
	return in == ExitCommand && breakline
}
