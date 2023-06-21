package config

import (
	"os"

	dallecfg "github.com/StepanTita/go-BingDALLE/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config interface {
	Logger
	Runtime
	Networker
	Prompter
	dallecfg.Authenticator
}

type config struct {
	Logger
	Runtime
	Networker
	Prompter
	dallecfg.Authenticator
}

type CliConfig struct {
	LogLevel string
	Rich     bool

	Proxy   string
	WssLink string

	Style   string
	Prompt  string
	Context string
	Locale  string

	DaLLe struct {
		ApiURL  string
		UCookie string
	}
}

type YamlGPTConfig struct {
	LogLevel string `yaml:"log_level"`
	Rich     bool   `yaml:"rich"`

	Proxy   string `yaml:"proxy"`
	WssLink string `yaml:"wss_link"`

	Style   string `yaml:"style"`
	Prompt  string `yaml:"prompt"`
	Context string `yaml:"context"`
	Locale  string `yaml:"locale"`

	DaLLe dallecfg.YamlDALLEConfig `yaml:"dalle"`
}

type yamlConfig struct {
	GPT YamlGPTConfig `yaml:"gpt"`
}

func NewFromGPTConfig(cfg YamlGPTConfig) Config {
	cfg.DaLLe.LogLevel = cfg.LogLevel
	cfg.DaLLe.Proxy = cfg.Proxy

	return &config{
		Logger:        NewLogger(cfg.LogLevel),
		Runtime:       NewRuntime(Version, cfg.Rich),
		Networker:     NewNetworker(cfg.DaLLe.ApiUrl, cfg.Proxy, cfg.WssLink),
		Prompter:      NewPrompter(cfg.Style, cfg.Prompt, cfg.Context, cfg.Locale),
		Authenticator: dallecfg.NewAuthenticator(cfg.DaLLe.UCookie),
	}
}

func NewFromFile(path string) Config {
	cfg := yamlConfig{}

	yamlConfig, err := os.ReadFile(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read config %s", path))
	}

	err = yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		panic(errors.Wrapf(err, "failed to unmarshal config %s", path))
	}

	return NewFromGPTConfig(cfg.GPT)
}

func NewFromCLI(cfg CliConfig) Config {
	return &config{
		Logger:    NewLogger(cfg.LogLevel),
		Runtime:   NewRuntime(Version, cfg.Rich),
		Networker: NewNetworker(cfg.DaLLe.ApiURL, cfg.Proxy, cfg.WssLink),
		Prompter:  NewPrompter(cfg.Style, cfg.Prompt, cfg.Context, cfg.Locale),

		Authenticator: dallecfg.NewAuthenticator(cfg.DaLLe.UCookie),
	}
}
