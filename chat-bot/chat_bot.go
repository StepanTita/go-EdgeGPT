package chat_bot

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common"
	"github.com/StepanTita/go-EdgeGPT/common/config"
	"github.com/StepanTita/go-EdgeGPT/common/convert"
	chat_hub "github.com/StepanTita/go-EdgeGPT/internal/services/chat-hub"
	"github.com/StepanTita/go-EdgeGPT/internal/services/conversation"
)

type ChatBot interface {
	Init(ctx context.Context, conversationStyle string, options ...string) error
	Ask(ctx context.Context, prompt, conversationStyle string, searchResult bool, options ...string) (<-chan ParsedFrame, error)
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

func (c *chatBot) Init(ctx context.Context, conversationStyle string, options ...string) error {
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

	msgsChan, err := c.chatHub.AskStream(ctx, c.cfg.InitialPrompt(), conversationStyle, false, options...)
	if err != nil {
		return errors.Wrap(err, "failed to ask stream")
	}
	// just ignoring the reply to initial prompt
	for range msgsChan {
	}
	return nil
}

func (c *chatBot) Ask(ctx context.Context, prompt, conversationStyle string, searchResult bool, options ...string) (<-chan ParsedFrame, error) {
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
		var respTxt, adaptiveCardsTxt string
		suggestedResponses := make([]string, 0, 10)

		var currMessageType *string = nil

		for msg := range msgsChan {
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

			out <- ParsedFrame{
				Text:               respTxt,
				AdaptiveCards:      adaptiveCardsTxt,
				Links:              common.ExtractURLs(adaptiveCardsTxt),
				Wrap:               wrap,
				Skip:               skip,
				SuggestedResponses: suggestedResponses,
			}
		}

		defer close(out)
	}()

	return out, nil
}
