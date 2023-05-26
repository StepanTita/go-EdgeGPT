package communicator

import (
	"context"
	"fmt"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/buger/goterm"
	"github.com/c-bata/go-prompt"
	"github.com/sirupsen/logrus"

	"github.com/briandowns/spinner"

	chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"
	"github.com/StepanTita/go-EdgeGPT/common/config"
	"github.com/StepanTita/go-EdgeGPT/common/terminal"
)

type Communicator struct {
	log *logrus.Entry

	cfg config.Config

	bot chat_bot.ChatBot

	suggestions []string
}

func New(cfg config.Config) *Communicator {
	return &Communicator{
		log: cfg.Logging().WithField("service", "[COMMUNICATOR]"),
		cfg: cfg,
		bot: chat_bot.New(cfg),
	}
}

func (c *Communicator) Run(ctx context.Context) error {
	inputSession := prompt.New(c.executorWithContext(ctx), c.completer,
		// options
		prompt.OptionSetExitCheckerOnInput(checkExit),
		prompt.OptionPrefix("User >>> "),
	)

	inputSession.Run()
	return nil
}

func (c *Communicator) executorWithContext(ctx context.Context) func(t string) {
	return func(t string) {
		area := terminal.NewArea()

		prefix := "Edge-GPT >>> "

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond,
			spinner.WithSuffix(fmt.Sprintf(" %s", prefix)),
			spinner.WithHiddenCursor(false),
		)

		s.Start()

		parsedResponsesChan, err := c.bot.Ask(ctx, t, c.cfg.Style(), true)
		if err != nil {
			c.log.WithError(err).Error("failed to ask bot")
			return
		}

		currText := prefix
		text := ""

		final := false
		for resp := range parsedResponsesChan {
			c.suggestions = resp.SuggestedResponses

			if resp.Skip {
				continue
			}

			text = resp.Text
			if c.cfg.Rich() {
				w := goterm.Width()
				text = string(markdown.Render(text, w, 0))
			}

			if !strings.HasSuffix(currText, prefix) {
				currText += prefix
			}

			s.Lock()
			if resp.Wrap {
				currText += text
				if !strings.HasSuffix(currText, "\n") {
					currText += "\n"
				}

				area.Update(currText)

				final = true
			} else {
				area.Update(currText + text)

				final = false
			}
			s.Unlock()
		}

		s.Stop()

		if final {
			area.Update(currText)
		} else {
			area.Update(currText + text)
		}
		goterm.Flush()
	}
}

func (c *Communicator) completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: HelpCommand},
		{Text: ExitCommand},
		{Text: ResetCommand},
	}

	for _, suggestion := range c.suggestions {
		s = append(s, prompt.Suggest{Text: suggestion})
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
