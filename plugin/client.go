package plugin

import (
	"fmt"
	"net/http"
)

// client provides utility for http request
type client struct {
	userAgent string
}

const defaultUserAgent = "mkr-plugin-installer/0.0.0"

// Get response from `url`
func (c *client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.getUA())
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

func (c *client) getUA() string {
	if c.userAgent != "" {
		return c.userAgent
	}
	return defaultUserAgent
}

func closeResponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
