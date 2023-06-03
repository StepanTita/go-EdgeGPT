package config

type Prompter interface {
	Style() string
	InitialPrompt() string
	Context() string
}

type prompter struct {
	style         string
	initialPrompt string
	context       string
}

func NewPrompter(style, prompt, context string) Prompter {
	return &prompter{
		style:         style,
		initialPrompt: prompt,
		context:       context,
	}
}

func (p prompter) Style() string {
	return p.style
}

func (p prompter) InitialPrompt() string {
	return p.initialPrompt
}

func (p prompter) Context() string {
	if p.context == "" {
		return defaultContext
	}
	return p.context
}
