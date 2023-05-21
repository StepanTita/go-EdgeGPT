package config

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/StepanTita/go-EdgeGPT/common/convert"
)

type Networker interface {
	Proxy() *url.URL
	WssLink() string
	Cookies() []map[string]any
}

type networker struct {
	proxy      string
	wssLink    string
	cookieFile string
}

func NewNetworker(cliConfig CliConfig) Networker {
	return &networker{
		proxy:      cliConfig.Proxy,
		wssLink:    cliConfig.WssLink,
		cookieFile: cliConfig.CookieFile,
	}
}

func getEnvFirstNotEmptyOrNil(names ...string) *string {
	for _, name := range names {
		if os.Getenv(name) != "" {
			return &name
		}
	}
	return nil
}

func (n networker) Proxy() *url.URL {
	if n.proxy == "" {
		proxy := getEnvFirstNotEmptyOrNil("all_proxy", "ALL_PROXY", "https_proxy", "HTTPS_PROXY")
		if proxy == nil {
			return nil
		}
		n.proxy = convert.FromPtr(proxy)
	}

	u, err := url.Parse(n.proxy)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse proxy url: %s", n.proxy))
	}

	// TODO: remove when http.Client would support socks5h
	if u.Scheme == "socks5h" {
		u.Scheme = "socks5"
	}
	return u
}

func (n networker) WssLink() string {
	return n.wssLink
}

func (n networker) Cookies() []map[string]any {
	f, err := os.Open(n.cookieFile)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open cookies file. Please make sure that specified path is valid: %s", n.cookieFile))
	}

	var rawCookies []map[string]any
	if err := json.NewDecoder(f).Decode(&rawCookies); err != nil {
		panic(errors.Wrap(err, "failed to decode cookies"))
	}
	return rawCookies
}
