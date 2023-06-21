package communicator

import (
	"context"
	"os"
	"runtime/debug"
	"time"

	"github.com/c-bata/go-prompt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/config"
)

type Communicator struct {
	log *logrus.Entry

	cfg config.Config

	renderer *renderer
}

func New(cfg config.Config) *Communicator {
	r := lipgloss.NewRenderer(os.Stderr, termenv.WithColorCache(true))

	opts := []tea.ProgramOption{tea.WithOutput(r.Output())}

	if !isatty.IsTerminal(os.Stdin.Fd()) {
		opts = append(opts, tea.WithInput(nil))
	}

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
		c.renderer = c.renderer.withContent("")

		if !c.checkCommand(t) {
			return
		}

		deadlineCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
		defer cancel()

		if err := c.renderer.run(deadlineCtx); err != nil {
			c.log.WithError(err).Error("failed to run renderer")
		}

		defer func() {
			if rvr := recover(); rvr != nil {
				c.log.Error("communicator panicked: ", rvr, string(debug.Stack()))
				c.log.Info("bot communication will be reset")
				c.reset()
			}
		}()
	}
}

func (c *Communicator) completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		// TODO: help command
		//{Text: HelpCommand},
		{Text: ExitCommand},
		{Text: ResetCommand},
	}

	for _, suggestion := range c.renderer.getSuggestions() {
		s = append(s, prompt.Suggest{Text: suggestion})
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (c *Communicator) checkCommand(t string) (exit bool) {
	switch t {
	case ResetCommand:
		c.reset()
		return true
	case ExitCommand:
		return false
	}
	return true
}

func (c *Communicator) reset() {
	c.renderer = c.renderer.withState(startState)
	c.renderer.withContent("")
}
