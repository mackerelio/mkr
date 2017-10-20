package plugin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestClientGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/200" {
			fmt.Fprint(w, "200 OK from "+req.Header.Get("User-Agent"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := &client{}

	{
		// Response is 200
		resp, err := c.get(ts.URL + "/200")
		defer resp.Body.Close()

		assert.NoError(t, err, "get finished successfully")

		b, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "200 OK from mkr-plugin-installer/0.0.0", string(b), "Getting response is success, and request has valid User-Agent")
	}

	{
		// Response is 404
		resp, err := c.get(ts.URL + "/404")

		assert.Nil(t, resp, "Return nothing as resp")
		assert.Error(t, err, "get failed")
		assert.Equal(t,
			fmt.Sprintf("http response not OK. code: 404, url: %s/404", ts.URL),
			err.Error(),
			"err has correct message",
		)
	}
}

func TestClientGetUA(t *testing.T) {
	c := &client{}
	assert.Equal(t, "mkr-plugin-installer/0.0.0", c.getUA(), "Returns default user agent")

	c = &client{userAgent: "custom-user-agent"}
	assert.Equal(t, "custom-user-agent", c.getUA(), "Returns custom user agent")
}
