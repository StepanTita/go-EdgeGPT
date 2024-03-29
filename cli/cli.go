package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/StepanTita/go-EdgeGPT/config"
	"github.com/StepanTita/go-EdgeGPT/services/communicator"
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
				Name:        "rich",
				Usage:       "Apply markdown renderer to output",
				Value:       false,
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
				Required:    false,
				Destination: &cliConfig.Prompt,
			},
			&cli.StringFlag{
				Name:        "context",
				Usage:       "Bot context to include in every request",
				Category:    "Bot:",
				Destination: &cliConfig.Context,
			},
			&cli.StringFlag{
				Name:        "locale",
				Usage:       "Locale for bot to use",
				Category:    "Bot:",
				Value:       "en",
				Destination: &cliConfig.Locale,
			},
			// DaLLe
			&cli.StringFlag{
				Name:        "bing-url",
				Usage:       "Bing URL (e.g. https://www.bing.com)",
				Category:    "DaLLe:",
				Required:    false,
				Value:       "https://www.bing.com",
				Destination: &cliConfig.DaLLe.ApiURL,
			},
			&cli.StringFlag{
				Name:        "u-auth-cookie",
				Usage:       "Cookie value to authenticate request",
				Category:    "DaLLe:",
				Required:    false,
				Destination: &cliConfig.DaLLe.UCookie,
			},
		},
		Commands: cli.Commands{
			{
				Name:  "run",
				Usage: "run EdgeGPT daemon",
				Action: func(c *cli.Context) error {
					cfg := config.NewFromCLI(cliConfig)
					log := cfg.Logging()

					fmt.Print("\033[H\033[2J")

					defer func() {
						if rvr := recover(); rvr != nil {
							log.Error("internal panicked: ", rvr, string(debug.Stack()))
						}
					}()

					comm := communicator.New(cfg)

					return comm.Run(c.Context)
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
