package config

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type Prompter interface {
	Style() string
	InitialPrompt() string
	Context() string
	Locale() string
	Language() string
}

type prompter struct {
	style         string
	initialPrompt string
	context       string
	locale        string
}

func NewPrompter(style, prompt, context, locale string) Prompter {
	return &prompter{
		style:         style,
		initialPrompt: prompt,
		context:       context,
		locale:        locale,
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

func (p prompter) Locale() string {
	return p.locale
}

func (p prompter) Language() string {
	return display.English.Tags().Name(language.Make(p.locale))
}
