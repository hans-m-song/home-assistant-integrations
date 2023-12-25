package util

import (
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
