package plugin

import (
	"fmt"
	"net/http"
)

// client provides utilities for http request
type client struct{}

const userAgent = "mkr-plugin-installer/0.0.0"

// Get response from `url`
func (c *client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http response not OK. code: %d, url: %s", resp.StatusCode, url)
		return nil, err
	}

	return resp, nil
}
