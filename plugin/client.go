package plugin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// client provides utilities for http request
type client struct{}

const userAgent = "mkr-plugin-installer/0.0.0"

// Get response from `url`
func (c *client) get(ctx context.Context, url string) (*http.Response, error) {
	resp, err := func() (*http.Response, error) {
		if strings.HasPrefix(url, "file:///") {
			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			return (&http.Client{Transport: t}).Get(url)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)
		return http.DefaultClient.Do(req)
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("http response not OK. code: %d, url: %s", resp.StatusCode, url)
	}
	return resp, nil
}
