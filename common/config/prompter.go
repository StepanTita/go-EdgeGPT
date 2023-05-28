package config

type Prompter interface {
	Style() string
	InitialPrompt() string
	Context() string
	AdaptiveCards() bool
}

type prompter struct {
	style         string
	initialPrompt string
	context       string
	adaptiveCards bool
}

func NewPrompter(cliConfig CliConfig) Prompter {
	return &prompter{
		style:         cliConfig.Style,
		initialPrompt: cliConfig.Prompt,
		context:       cliConfig.Context,
		adaptiveCards: cliConfig.AdaptiveCards,
	}
}

func (p prompter) Style() string {
	return p.style
}

func (p prompter) InitialPrompt() string {
	return p.initialPrompt
}

func (p prompter) AdaptiveCards() bool {
	return p.adaptiveCards
}

func (p prompter) Context() string {
	if p.context == "" {
		return defaultContext
	}
	return p.context
}
