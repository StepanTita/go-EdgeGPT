package cli

import (
	"fmt"
	"os"

	"atomicgo.dev/cursor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"
	"github.com/StepanTita/go-EdgeGPT/common/config"
)

func Run(args []string) bool {
	var cliConfig config.CliConfig

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "Log level ['debug', 'info', 'warn', 'error', 'fatal']",
				Value:       "info",
				Category:    "Miscellaneous:",
				Destination: &cliConfig.LogLevel,
			},
			&cli.BoolFlag{
				Name:        "no-stream",
				Usage:       "Do not stream and format the output",
				Value:       true,
				Category:    "Miscellaneous:",
				Destination: &cliConfig.NoStream,
			},
			&cli.BoolFlag{
				Name:        "rich",
				Usage:       "Apply markdown renderer to output",
				Value:       true,
				Category:    "Miscellaneous:",
				Destination: &cliConfig.Rich,
			},
			&cli.StringFlag{
				Name:        "proxy",
				Usage:       "Proxy URL (e.g. socks5://127.0.0.1:1080)",
				Category:    "Networking:",
				Required:    false,
				Destination: &cliConfig.Proxy,
			},
			&cli.StringFlag{
				Name:        "wss-link",
				Usage:       "WSS URL(e.g. wss://sydney.bing.com/sydney/ChatHub)",
				Value:       "wss://sydney.bing.com/sydney/ChatHub",
				Category:    "Networking:",
				Destination: &cliConfig.WssLink,
			},
			&cli.StringFlag{
				Name:        "cookie-file",
				Usage:       "Cookie file used for authentication (defaults to COOKIE_FILE environment variable)",
				Value:       os.Getenv("COOKIE_FILE"),
				Category:    "Networking:",
				Required:    false,
				Destination: &cliConfig.CookieFile,
			},
			&cli.StringFlag{
				Name:        "style",
				Usage:       "Style of the conversation with bot ['creative', 'balanced', 'precise']",
				Value:       "balanced",
				Category:    "Bot:",
				Destination: &cliConfig.Style,
			},
			&cli.StringFlag{
				Name:        "prompt",
				Usage:       "Prompt to start with",
				Category:    "Bot:",
				Destination: &cliConfig.Prompt,
			},
		},
		Commands: cli.Commands{
			{
				Name:  "run",
				Usage: "run EdgeGPT daemon",
				Action: func(c *cli.Context) error {
					cfg := config.New(cliConfig)
					log := cfg.Logging()

					defer func() {
						if rvr := recover(); rvr != nil {
							log.Error("internal panicked: ", rvr)
						}
					}()

					bot := chat_bot.New(cfg)

					out, err := bot.Ask(c.Context, cfg.InitialPrompt(), "creative", false)
					if err != nil {
						log.WithError(err).Error("failed to ask bot")
					}

					for resp := range out {
						cursor.StartOfLine()
						cursor.ClearLine()
						fmt.Print(resp.Text)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(args); err != nil {
		logrus.Error(err, ": service initialization failed")
		return false
	}

	return true
}
