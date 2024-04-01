package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

func HTTPGet(ctx context.Context, url string, query url.Values) (int, []byte, error) {
	if query != nil {
		url += "?" + query.Encode()
	}

	log.Trace().Str("url", url).Str("method", http.MethodGet).Msg("sending request")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to do request: %s", err)
	}

	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read response body: %s", err)
	}

	log.Trace().Str("url", url).Str("method", http.MethodGet).Bytes("body", raw).Msg("received response")
	return resp.StatusCode, raw, nil
}

type LogRoundTripper struct {
	Name      string
	Transport http.RoundTripper
}

func (t *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log := log.With().
		Str("transport", t.Name).
		Str("url", req.URL.String()).
		Str("method", req.Method).
		Logger()

	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		log = log.With().AnErr("request_err", err).Logger()
	}

	if resp == nil {
		log.Error().Msg("response is nil")
		return resp, err
	}

	raw, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	log = log.With().AnErr("body_err", err).Str("body", string(raw)).Logger()

	resp.Body = io.NopCloser(bytes.NewReader(raw))
	log.Trace().Send()

	return resp, nil
}
