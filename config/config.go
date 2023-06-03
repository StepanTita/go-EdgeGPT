package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

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
	Rich     bool

	Proxy   string
	WssLink string

	Style   string
	Prompt  string
	Context string
}

type YamlGPTConfig struct {
	LogLevel string `yaml:"log_level"`
	Rich     bool   `yaml:"rich"`

	Proxy   string `yaml:"proxy"`
	WssLink string `yaml:"wss_link"`

	Style   string `yaml:"style"`
	Prompt  string `yaml:"prompt"`
	Context string `yaml:"context"`
}

type yamlConfig struct {
	GPT YamlGPTConfig `yaml:"gpt"`
}

func NewFromGPTConfig(cfg YamlGPTConfig) Config {
	return &config{
		Logger:    NewLogger(cfg.LogLevel),
		Runtime:   NewRuntime(Version, cfg.Rich),
		Networker: NewNetworker(cfg.Proxy, cfg.WssLink),
		Prompter:  NewPrompter(cfg.Style, cfg.Prompt, cfg.Context),
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
		Networker: NewNetworker(cfg.Proxy, cfg.WssLink),
		Prompter:  NewPrompter(cfg.Style, cfg.Prompt, cfg.Context),
	}
}
