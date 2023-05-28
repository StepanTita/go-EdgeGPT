package communicator

import (
	"context"

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common/config"
)

type Communicator struct {
	log *logrus.Entry

	cfg config.Config

	renderer *renderer
}

func New(cfg config.Config, r *lipgloss.Renderer) *Communicator {
	return &Communicator{
		log: cfg.Logging().WithField("service", "[COMMUNICATOR]"),
		cfg: cfg,

		renderer: newRenderer(cfg, r),
	}
}

func (c *Communicator) Run(ctx context.Context) error {
	// Init prompt
	if err := c.renderer.run(ctx); err != nil {
		return errors.Wrap(err, "failed to run renderer")
	}

	inputSession := prompt.New(c.executorWithContext(ctx), c.completer,
		// options
		prompt.OptionSetExitCheckerOnInput(checkExit),
		prompt.OptionPrefix("> "),
	)

	inputSession.Run()
	return nil
}

func (c *Communicator) executorWithContext(ctx context.Context) func(t string) {
	return func(t string) {
		c.renderer = c.renderer.withState(completionState)
		c.renderer = c.renderer.withInput(t)

		c.checkCommand(t)

		if err := c.renderer.run(ctx); err != nil {
			c.log.WithError(err).Error("failed to run renderer")
		}
	}
}

func (c *Communicator) completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: HelpCommand},
		{Text: ExitCommand},
		{Text: ResetCommand},
	}

	for _, suggestion := range c.renderer.getSuggestions() {
		s = append(s, prompt.Suggest{Text: suggestion})
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (c *Communicator) checkCommand(t string) {
	if t == ResetCommand {
		c.renderer = c.renderer.withState(startState)
		c.renderer.withContent("")
	}
}
