package chat_bot

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common/config"
	"github.com/StepanTita/go-EdgeGPT/common/convert"
	chat_hub "github.com/StepanTita/go-EdgeGPT/internal/services/chat-hub"
	"github.com/StepanTita/go-EdgeGPT/internal/services/conversation"
)

type ChatBot struct {
	log *logrus.Entry

	cfg config.Config
}

func New(cfg config.Config) ChatBot {
	return ChatBot{
		log: cfg.Logging().WithField("service", "[CHAT-BOT]"),

		cfg: cfg,
	}
}

func (c ChatBot) Ask(ctx context.Context, prompt, conversationStyle string, searchResult bool) (<-chan ParsedFrame, error) {
	conv := conversation.New(c.cfg)

	state, err := conv.Create(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new conversation")
	}
	c.log.WithFields(logrus.Fields{
		"result":          state.Result,
		"conversation_id": state.ConversationID,
		"client_id":       state.ClientID,
	}).Info("Created conversation")

	chatHub := chat_hub.New(c.cfg, state)

	msgsChan, err := chatHub.AskStream(ctx, prompt, conversationStyle, searchResult)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ask stream")
	}

	out, err := c.processMessages(msgsChan)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process messages")
	}

	return out, nil
}

func (c ChatBot) processMessages(msgsChan <-chan chat_hub.ResponseMessage) (<-chan ParsedFrame, error) {
	out := make(chan ParsedFrame)

	go func() {
		var respTxt string

		for msg := range msgsChan {
			if msg.Type == 1 && msg.Arguments[0].Messages != nil {
				msgText := convert.FromPtr(msg.Arguments[0].Messages[0].AdaptiveCards[0].Body[0].Text)
				if strings.HasSuffix(respTxt, "\n") {
					continue
				}
				respTxt = msgText

				out <- ParsedFrame{
					Text: respTxt,
					Skip: false,
				}
			} else if msg.Type == 2 {
				if msg.Item.Result.Error != nil {
					c.log.WithFields(logrus.Fields{
						"value":   msg.Item.Result.Value,
						"message": msg.Item.Result.Message,
					}).Fatal("Some error occurred")
				}
			}
		}
	}()

	return out, nil
}
