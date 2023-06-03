package communicator

import (
	"context"
	"fmt"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/sirupsen/logrus"

	chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"
	"github.com/StepanTita/go-EdgeGPT/config"
	terminal2 "github.com/StepanTita/go-EdgeGPT/terminal"
)

type renderer struct {
	log *logrus.Entry
	ctx context.Context
	cfg config.Config

	bot chat_bot.ChatBot

	suggestions []string

	userInput string
	content   string

	parsedResponsesChan <-chan chat_bot.ParsedFrame

	// tea Model fields
	prefix string

	options []tea.ProgramOption

	program  *tea.Program
	Error    *communicatorError
	state    state
	retries  int
	styles   terminal2.Styles
	renderer *lipgloss.Renderer
	anim     tea.Model

	width  int
	height int
}

// chatInput is a tea.Msg that wraps the content read from stdin.
type chatInput struct {
	content       string
	adaptiveCards string
}

// chatOutput a tea.Msg that wraps the content returned from openai.
type chatOutput struct {
	content string
}

type initPrompt struct{}

func newRenderer(cfg config.Config, r *lipgloss.Renderer, opts ...tea.ProgramOption) *renderer {
	styles := terminal2.NewStyles(r)

	prefix := styles.Prefix.Render("Edge-GPT >")
	if r.ColorProfile() == termenv.TrueColor {
		prefix = terminal2.MakeGradientText(styles.Prefix, "Edge-GPT >")
	}

	rend := &renderer{
		log: cfg.Logging().WithField("service", "[RENDERER]"),
		cfg: cfg,

		bot: chat_bot.New(cfg),

		prefix: prefix,

		state:    startState,
		renderer: r,
		styles:   styles,

		options: opts,
	}

	rend.program = tea.NewProgram(rend, opts...)
	return rend
}

func (r *renderer) withContext(ctx context.Context) *renderer {
	r.ctx = ctx
	return r
}

func (r *renderer) withState(state state) *renderer {
	r.state = state
	return r
}

func (r *renderer) withInput(input string) *renderer {
	r.userInput = input
	return r
}

func (r *renderer) withContent(content string) *renderer {
	r.content = content
	return r
}

func (r *renderer) getSuggestions() []string {
	return r.suggestions
}

func (r *renderer) run(ctx context.Context) error {
	*r = *r.withContext(ctx)

	r.program = tea.NewProgram(r, r.options...)

	_, err := r.program.Run()

	return err
}

// Init implements tea.Model.
func (r *renderer) Init() tea.Cmd {
	if r.state == startState {
		return func() tea.Msg {
			return initPrompt{}
		}
	}

	var err error
	r.parsedResponsesChan, err = r.bot.Ask(r.ctx, r.userInput, r.cfg.Style(), true)
	if err != nil {
		r.log.WithError(err).Error("failed to ask bot")
		return func() tea.Msg {
			return communicatorError{
				err:    err,
				reason: "failed to ask bot",
			}
		}
	}
	return r.readFrame
}

func (r *renderer) initCall() tea.Msg {
	r.state = initRunningState

	if err := r.bot.Init(r.ctx); err != nil {
		r.log.WithError(err).Error("failed to initialize bot state")
		return communicatorError{
			err:    err,
			reason: "failed to initialize bot state",
		}
	}

	if r.cfg.InitialPrompt() != "" {
		if err := r.bot.InitPrompt(r.ctx, r.cfg.Style()); err != nil {
			r.log.WithError(err).Error("failed to run initial prompt")
			return communicatorError{
				err:    err,
				reason: "failed to run initial prompt",
			}
		}
	}

	r.state = initCompletedState
	return initPrompt{}
}

// Update implements tea.Model.
func (r *renderer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initPrompt:
		if r.state == initCompletedState {
			r.anim = terminal2.NewCyclingChars(cyclingChars, generationText, r.renderer, r.styles)
			return r, tea.Quit
		} else if r.state == initRunningState {
			return r, r.anim.Init()
		}

		r.anim = terminal2.NewCyclingChars(cyclingChars, initStatusText, r.renderer, r.styles)
		return r, tea.Batch(r.anim.Init(), r.initCall)
	case chatInput:
		r.content = msg.content
		if msg.content != "" {
			s := fmt.Sprintf("%s\n\n%s", msg.content, msg.adaptiveCards)
			r.content = fmt.Sprintf("%s %s", r.prefix, strings.TrimPrefix(s, "\n"))
		}
		r.content = fmt.Sprintf("%s %s", r.prefix, strings.TrimPrefix(msg.adaptiveCards, "\n"))

		return r, tea.Batch(r.anim.Init(), r.readFrame)
	case chatOutput:
		r.content = msg.content
		return r, tea.Quit
	case communicatorError:
		r.Error = &msg
		r.state = errorState
		return r, tea.Quit
	case tea.WindowSizeMsg:
		r.width, r.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+с":
			return r, tea.Quit
		}
	}
	if r.state == initRunningState || r.state == completionState {
		var cmd tea.Cmd
		r.anim, cmd = r.anim.Update(msg)
		return r, cmd
	}
	return r, nil
}

// View implements tea.Model.
func (r *renderer) View() string {
	switch r.state {
	case errorState:
		return r.ErrorView()
	case completionState, initRunningState:
		return fmt.Sprintf("%s\n\n%s", r.anim.View(), r.FormattedOutput())
	}
	return r.FormattedOutput()
}

// ErrorView renders the currently set modsError
func (r *renderer) ErrorView() string {
	const maxWidth = 120
	const horizontalPadding = 2
	w := r.width - (horizontalPadding * 2)
	if w > maxWidth {
		w = maxWidth
	}
	s := r.renderer.NewStyle().Width(w).Padding(0, horizontalPadding)
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n",
		s.Render(r.styles.ErrorHeader.String(), r.Error.reason),
		s.Render(r.styles.ErrorDetails.Render(r.Error.Error())),
	)
}

// FormattedOutput returns the response from OpenAI with the user configured
// prefix and standard in settings.
func (r *renderer) FormattedOutput() string {
	if r.cfg.Rich() {
		return string(markdown.Render(r.content, r.width, 0))
	}
	return r.content
}

// readFrame reads single frame from the input channel
func (r *renderer) readFrame() tea.Msg {
	parsedFrame, ok := <-r.parsedResponsesChan
	if !ok {
		r.state = completionDoneState
		return chatOutput{content: r.content}
	}
	r.suggestions = parsedFrame.SuggestedResponses
	if parsedFrame.Skip {
		return chatInput{content: r.content}
	}
	return chatInput{content: parsedFrame.Text, adaptiveCards: parsedFrame.AdaptiveCards}
}
