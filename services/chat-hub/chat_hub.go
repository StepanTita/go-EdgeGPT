package chat_hub

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/config"
	"github.com/StepanTita/go-EdgeGPT/services/conversation"
)

type ChatHub struct {
	log *logrus.Entry
	cfg config.Config

	request chatHubRequest

	dialer ws.Dialer

	conn net.Conn
}

func New(cfg config.Config, state *conversation.State) ChatHub {
	return ChatHub{
		log: cfg.Logging().WithField("service", "[CHAT-HUB]"),
		cfg: cfg,

		request: newChatHubRequest(state),

		dialer: ws.Dialer{},
	}
}

func (c *ChatHub) initialHandshake(ctx context.Context) error {
	conn, _, _, err := c.dialer.Dial(ctx, c.cfg.WssLink())
	if err != nil {
		return errors.Wrap(err, "failed to dial wss. Initial handshake failed")
	}

	c.conn = conn
	if err = wsutil.WriteClientMessage(c.conn, ws.OpText, appendIdentifier(json.RawMessage(`{"protocol": "json", "version": 1}`))); err != nil {
		return errors.Wrap(err, "failed to write an initial message to the buffer")
	}

	var out []byte
	if _, err = wsutil.ReadServerMessage(c.conn, []wsutil.Message{
		{
			OpCode:  ws.OpText,
			Payload: out,
		},
	}); err != nil {
		return errors.Wrap(err, "failed to read an initial message from the buffer")
	}
	return nil
}

func (c *ChatHub) AskStream(ctx context.Context, prompt string, conversationalStyle string, searchResult bool, options ...string) (<-chan ResponseMessage, error) {
	if err := c.initialHandshake(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to reestablish connection")
	}

	if err := c.request.Update(prompt, conversationalStyle, options, searchResult); err != nil {
		return nil, errors.Wrap(err, "failed to update chat hub request")
	}

	msg, err := c.request.EncodeJson()
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode json message")
	}

	if err := wsutil.WriteClientMessage(c.conn, ws.OpText, appendIdentifier(msg)); err != nil {
		return nil, errors.Wrap(err, "failed to write message")
	}

	out := make(chan ResponseMessage)

	go func() {
		defer close(out)
		for {
			rawMsg, _, err := wsutil.ReadServerData(c.conn)
			if err != nil {
				c.log.WithError(err).Error("failed to read message from the server")
				return
			}

			responseMsg := ResponseMessage{}
			for _, obj := range strings.Split(string(rawMsg), delimeter) {
				if obj == "" {
					continue
				}
				if err := json.Unmarshal([]byte(obj), &responseMsg); err != nil {
					c.log.WithError(err).Error("failed to unmarshall response message")
					return
				}

				out <- responseMsg

				if responseMsg.Type == 2 {
					c.conn.Close()
					return
				}
			}
		}
	}()

	return out, nil
}
