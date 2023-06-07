package chat_bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common"
	"github.com/StepanTita/go-EdgeGPT/common/convert"
	"github.com/StepanTita/go-EdgeGPT/config"
	"github.com/StepanTita/go-EdgeGPT/services/chat-hub"
	"github.com/StepanTita/go-EdgeGPT/services/conversation"
)

const (
	markdownPrefix = "Format the response as Markdown."
	languagePrefix = "In your response use: %s. Do not provide translation to english if language is not english."
)

type ChatBot interface {
	Init(ctx context.Context) error
	InitPrompt(ctx context.Context, conversationStyle string, options ...string) error
	Ask(ctx context.Context, prompt, context, conversationStyle string, searchResult bool, language string, options ...string) (<-chan ParsedFrame, error)
	EstimatePrompt(prompt, context, language string) int
}

type chatBot struct {
	log *logrus.Entry

	cfg config.Config

	conv    conversation.Conversation
	chatHub chat_hub.ChatHub

	state *conversation.State
}

func New(cfg config.Config) ChatBot {
	return &chatBot{
		log: cfg.Logging().WithField("service", "[CHAT-BOT]"),

		cfg: cfg,

		conv: conversation.New(cfg),
	}
}

func (c *chatBot) Init(ctx context.Context) error {
	var err error
	c.state, err = c.conv.Create(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create new conversation")
	}

	c.log.WithFields(logrus.Fields{
		"result":          c.state.Result,
		"conversation_id": c.state.ConversationID,
		"client_id":       c.state.ClientID,
	}).Info("Created conversation")

	c.chatHub = chat_hub.New(c.cfg, c.state)
	return nil
}

func (c *chatBot) InitPrompt(ctx context.Context, conversationStyle string, options ...string) error {
	msgsChan, err := c.chatHub.AskStream(ctx, c.cfg.InitialPrompt(), conversationStyle, false, options...)
	if err != nil {
		return errors.Wrap(err, "failed to ask stream")
	}
	// just ignoring the reply to initial prompt
	for range msgsChan {
	}
	return nil
}

func (c *chatBot) EstimatePrompt(prompt, context, language string) int {
	if context != "" {
		prompt = fmt.Sprintf("%s\n\n%s", context, prompt)
	}

	if c.cfg.Rich() {
		prompt = strings.ReplaceAll(prompt, "{{markdown}}", markdownPrefix)
	}

	prompt = strings.ReplaceAll(prompt, "{{language}}", fmt.Sprintf(languagePrefix, language))
	return len(prompt)
}

func (c *chatBot) Ask(ctx context.Context, prompt, context, conversationStyle string, searchResult bool, language string, options ...string) (<-chan ParsedFrame, error) {

	if context != "" {
		prompt = fmt.Sprintf("%s\n\n%s", context, prompt)
	}

	if c.cfg.Rich() {
		prompt = strings.ReplaceAll(prompt, "{{markdown}}", markdownPrefix)
	}

	prompt = strings.ReplaceAll(prompt, "{{language}}", fmt.Sprintf(languagePrefix, language))

	msgsChan, err := c.chatHub.AskStream(ctx, prompt, conversationStyle, searchResult, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ask stream")
	}

	parsedFramesChan, err := c.processMessages(ctx, msgsChan)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process messages")
	}

	return parsedFramesChan, nil
}

func (c *chatBot) processMessages(ctx context.Context, msgsChan <-chan chat_hub.ResponseMessage) (<-chan ParsedFrame, error) {
	out := make(chan ParsedFrame)

	go func() {
		defer close(out)

		var respTxt, adaptiveCardsTxt, referencesText string
		suggestedResponses := make([]string, 0, 10)
		var links []ResponseLink

		var currMessageType *string = nil

		frame := 0
		for {
			c.log.WithField("frame", frame).Debug("Parsing frame...")
			frame++
			select {
			case <-ctx.Done():
				c.log.WithTime(common.CurrentTimestamp()).Error("deadline exceeded")
				return
			case msg, ok := <-msgsChan:
				if !ok {
					return
				}
				skip := false
				wrap := false
				respTxt = ""

				if msg.Type == 1 && msg.Arguments[0].Messages != nil {
					if msg.Arguments[0].Messages[0].MessageType == nil {
						adaptiveCardsTxt = convert.FromPtr(msg.Arguments[0].Messages[0].AdaptiveCards[0].Body[0].Text)
					} else if convert.FromPtr(msg.Arguments[0].Messages[0].MessageType) == MessageTypeDisengaged {
						respTxt = "The conversation has been stopped prematurely... Sorry, please, restart the conversation"
					} else {
						respTxt = convert.FromPtr(msg.Arguments[0].Messages[0].Text)
					}

					if msg.Arguments[0].Messages[0].MessageType != currMessageType {
						currMessageType = msg.Arguments[0].Messages[0].MessageType
						wrap = true
					}

					if convert.FromPtr(msg.Arguments[0].Messages[0].MessageType) == MessageTypeRenderCardRequest {
						skip = true
					}

				} else if msg.Type == 2 {
					if msg.Item.Result.Error != nil {
						c.log.WithFields(logrus.Fields{
							"value":   convert.FromPtr(msg.Item.Result.Value),
							"message": convert.FromPtr(msg.Item.Result.Message),
						}).Fatal("Some error occurred")
					}

					if len(msg.Item.Messages) > 0 {
						for _, item := range msg.Item.Messages {
							if convert.FromPtr(item.ContentOrigin) == "Apology" {
								adaptiveCardsTxt = convert.FromPtr(item.AdaptiveCards[0].Body[0].Text)
							}
						}

						lastId := len(msg.Item.Messages) - 1

						for _, suggestion := range msg.Item.Messages[lastId].SuggestedResponses {
							suggestedResponses = append(suggestedResponses, convert.FromPtr(suggestion.Text))
						}
					}
				}

				if respTxt == "" && adaptiveCardsTxt == "" {
					skip = true
				}

				if len(links) == 0 {
					referencesText, links = ExtractURLs(adaptiveCardsTxt)
				}

				out <- ParsedFrame{
					Text:               respTxt,
					AdaptiveCards:      strings.ReplaceAll(adaptiveCardsTxt, referencesText, ""),
					Links:              links,
					Wrap:               wrap,
					Skip:               skip,
					SuggestedResponses: suggestedResponses,
				}
			}
		}
	}()

	return out, nil
}
