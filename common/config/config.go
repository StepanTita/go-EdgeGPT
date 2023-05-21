package config

type Config interface {
	Logger
	Runtime
	Networker
	Prompter
}

type config struct {
	Logger
	Runtime
	Networker
	Prompter
}

type CliConfig struct {
	LogLevel string
	NoStream bool
	Rich     bool

	Proxy      string
	WssLink    string
	CookieFile string

	Style  string
	Prompt string
}

func New(cliConfig CliConfig) Config {
	return &config{
		Logger:    NewLogger(cliConfig.LogLevel),
		Runtime:   NewRuntime(Version, cliConfig),
		Networker: NewNetworker(cliConfig),
		Prompter:  NewPrompter(cliConfig),
	}
}
