package communicator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"atomicgo.dev/cursor"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"

	chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"
	"github.com/StepanTita/go-EdgeGPT/common/config"
)

type Communicator struct {
	log *logrus.Entry

	cfg config.Config

	bot *chat_bot.ChatBot

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
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond,
			spinner.WithSuffix(" Edge-GPT >>> "),
			spinner.WithHiddenCursor(false),
		)

		s.Start()

		parsedResponsesChan, err := c.bot.Ask(ctx, t, c.cfg.Style(), true)
		if err != nil {
			c.log.WithError(err).Error("failed to ask bot")
		}

		for resp := range parsedResponsesChan {
			c.suggestions = resp.SuggestedResponses

			text := resp.Text
			if c.cfg.Rich() {
				w, _, _ := term.GetSize(0)
				text = string(markdown.Render(text, w, 0))
			}

			if resp.Skip {
				continue
			}

			if resp.Wrap {
				text = strings.TrimSuffix(text, "\n")

				cursor.StartOfLine()
				cursor.ClearLine()
				fmt.Println(fmt.Sprintf("Edge-GPT >>> %s", text))
				s.Suffix = " Edge-GPT >>> "
			} else {
				s.Suffix = fmt.Sprintf(" Edge-GPT >>> %s", text)
			}
		}
		s.Stop()
		fmt.Print(s.Suffix[1:])
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
