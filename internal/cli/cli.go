package cli

import (
	"os"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/StepanTita/go-EdgeGPT/common/config"
	"github.com/StepanTita/go-EdgeGPT/internal/services/communicator"
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
			&cli.StringFlag{
				Name:        "context",
				Usage:       "Bot context to include in every request",
				Category:    "Bot:",
				Destination: &cliConfig.Context,
			},
		},
		Commands: cli.Commands{
			{
				Name:  "run",
				Usage: "run EdgeGPT daemon",
				Action: func(c *cli.Context) error {
					cfg := config.NewFromCLI(cliConfig)
					log := cfg.Logging()

					defer func() {
						if rvr := recover(); rvr != nil {
							log.Error("internal panicked: ", rvr, string(debug.Stack()))
						}
					}()

					renderer := lipgloss.NewRenderer(os.Stderr, termenv.WithColorCache(true))
					opts := []tea.ProgramOption{tea.WithOutput(renderer.Output())}
					if !isatty.IsTerminal(os.Stdin.Fd()) {
						opts = append(opts, tea.WithInput(nil))
					}

					comm := communicator.New(cfg, renderer)

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
