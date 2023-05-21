package connector

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/StepanTita/go-EdgeGPT/common"
	"github.com/StepanTita/go-EdgeGPT/common/config"
)

type Connector interface {
	Request(ctx context.Context, r RequestParams) (io.ReadCloser, int, error)
}

type connector struct {
	log *logrus.Entry

	cfg config.Config

	client http.Client
}

func New(cfg config.Config) Connector {
	return &connector{
		log: cfg.Logging().WithField("service", "[CONN]"),

		cfg: cfg,
		client: http.Client{
			Timeout: 0,
		},
	}
}

// Request TODO: add retry policy
func (c connector) Request(ctx context.Context, r RequestParams) (io.ReadCloser, int, error) {
	c.log.Debugf("Requesting, %s%s...", r.Url, r.Path)

	c.client.Transport = &http.Transport{Proxy: http.ProxyURL(c.cfg.Proxy())}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", r.Url, r.Path), nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create new request")
	}

	for _, cookie := range c.cfg.Cookies() {
		req.AddCookie(&http.Cookie{
			Name:  cookie["name"].(string),
			Value: cookie["value"].(string),
		})
	}

	req.Header = common.HEADERS_HTTPS

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to do client request")
	}
	c.log.WithField("status_code", resp.StatusCode).Debug("request completed with the status code")

	return resp.Body, resp.StatusCode, nil
}