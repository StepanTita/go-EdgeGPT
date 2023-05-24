package chat_bot

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common/config"
	"github.com/StepanTita/go-EdgeGPT/common/convert"
	chat_hub "github.com/StepanTita/go-EdgeGPT/internal/services/chat-hub"
	"github.com/StepanTita/go-EdgeGPT/internal/services/conversation"
)

type ChatBot interface {
	Ask(ctx context.Context, prompt, conversationStyle string, searchResult bool, options ...string) (<-chan ParsedFrame, error)
}

type chatBot struct {
	log *logrus.Entry

	cfg config.Config

	conv    conversation.Conversation
	chatHub chat_hub.ChatHub

	state *conversation.State
	once  *sync.Once
}

func New(cfg config.Config) ChatBot {
	return &chatBot{
		log: cfg.Logging().WithField("service", "[CHAT-BOT]"),

		cfg: cfg,

		conv: conversation.New(cfg),
		once: &sync.Once{},
	}
}

func (c *chatBot) Ask(ctx context.Context, prompt, conversationStyle string, searchResult bool, options ...string) (<-chan ParsedFrame, error) {
	c.once.Do(func() {
		var err error
		c.state, err = c.conv.Create(ctx)
		if err != nil {
			c.log.WithError(err).Fatal("failed to create new conversation")
		}

		c.log.WithFields(logrus.Fields{
			"result":          c.state.Result,
			"conversation_id": c.state.ConversationID,
			"client_id":       c.state.ClientID,
		}).Info("Created conversation")

		c.chatHub = chat_hub.New(c.cfg, c.state)

		msgsChan, err := c.chatHub.AskStream(ctx, c.cfg.InitialPrompt(), conversationStyle, searchResult, options...)
		// just ignoring the reply to initial prompt
		for range msgsChan {
		}
	})

	msgsChan, err := c.chatHub.AskStream(ctx, prompt, conversationStyle, searchResult, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ask stream")
	}

	parsedFramesChan, err := c.processMessages(msgsChan)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process messages")
	}

	return parsedFramesChan, nil
}

func (c *chatBot) processMessages(msgsChan <-chan chat_hub.ResponseMessage) (<-chan ParsedFrame, error) {
	out := make(chan ParsedFrame)

	go func() {
		var respTxt, updatedTxt string
		suggestedResponses := make([]string, 0, 10)

		for msg := range msgsChan {
			skip := false
			wrap := false
			respTxt = updatedTxt
			if msg.Type == 1 && msg.Arguments[0].Messages != nil && msg.Arguments[0].Messages[0].Text != nil {
				if c.cfg.AdaptiveCards() && msg.Arguments[0].Messages[0].MessageType == nil {
					respTxt = convert.FromPtr(msg.Arguments[0].Messages[0].AdaptiveCards[0].Body[0].Text)
				} else {
					respTxt = convert.FromPtr(msg.Arguments[0].Messages[0].Text)
				}
				updatedTxt = respTxt

				if msg.Arguments[0].Messages[0].MessageType != nil {
					updatedTxt = ""
					wrap = true
				}

				if convert.FromPtr(msg.Arguments[0].Messages[0].MessageType) == "RenderCardRequest" {
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
					lastId := len(msg.Item.Messages) - 1

					for _, suggestion := range msg.Item.Messages[lastId].SuggestedResponses {
						suggestedResponses = append(suggestedResponses, convert.FromPtr(suggestion.Text))
					}
				}
			}

			if respTxt == "" {
				skip = true
			}

			out <- ParsedFrame{
				Text:               respTxt,
				Wrap:               wrap,
				Skip:               skip,
				SuggestedResponses: suggestedResponses,
			}
		}

		defer close(out)
	}()

	return out, nil
}
