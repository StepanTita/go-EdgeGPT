package conversation

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common/config"
	"github.com/StepanTita/go-EdgeGPT/internal/services/connector"
)

type Conversation struct {
	log *logrus.Entry
	cfg config.Config

	conn connector.Connector
}

func New(cfg config.Config) Conversation {
	return Conversation{
		log:  cfg.Logging().WithField("service", "[CONVERSATION]"),
		cfg:  cfg,
		conn: connector.New(cfg),
	}
}

func (c Conversation) Create(ctx context.Context) (*State, error) {
	bodyReader, status, err := c.conn.Request(ctx, connector.RequestParams{
		// https://edge.churchless.tech TODO: might need to try that as well if the first one failed
		Url:  "https://edgeservices.bing.com",
		Path: "/edgesvc/turing/conversation/create",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send create conversation request")
	}

	if status != http.StatusOK {
		var body map[string]any
		if err := json.NewDecoder(bodyReader).Decode(&body); err != nil {
			return nil, errors.Wrapf(err, "failed to decode body with status code: %d", status)
		}
		c.log.WithFields(logrus.Fields{
			"body":   body,
			"status": status,
		}).Error("create conversation request failed")
		return nil, errors.New("")
	}

	state := State{}
	if err := json.NewDecoder(bodyReader).Decode(&state); err != nil {
		return nil, errors.Wrap(err, "failed to decode state")
	}
	return &state, nil
}
