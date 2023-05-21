package config

type Prompter interface {
	Style() string
	InitialPrompt() string
}

type prompter struct {
	style         string
	initialPrompt string
}

func NewPrompter(cliConfig CliConfig) Prompter {
	return &prompter{
		style:         cliConfig.Style,
		initialPrompt: cliConfig.Prompt,
	}
}

func (p prompter) Style() string {
	return p.style
}

func (p prompter) InitialPrompt() string {
	return p.initialPrompt
}
