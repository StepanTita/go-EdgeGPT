package communicator

import (
	"github.com/c-bata/go-prompt"

	"github.com/StepanTita/go-EdgeGPT/common/config"
)

type Communicator struct {
	cfg          config.Config
	inputSession *prompt.Prompt
}

func New(cfg config.Config) Communicator {
	return Communicator{
		cfg: cfg,
		inputSession: prompt.New(executor, completer,
			// options
			prompt.OptionSetExitCheckerOnInput(checkExit),
		),
	}
}
