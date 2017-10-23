package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Makes related environments empty,
	// and reset these after test
	origTokenEnv := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", origTokenEnv)
	os.Unsetenv("GITHUB_TOKEN")
	origGitConfigEnv := os.Getenv("GIT_CONFIG")
	defer os.Setenv("GIT_CONFIG", origGitConfigEnv)
	os.Unsetenv("GIT_CONFIG")
	os.Exit(m.Run())
}

func TestGetGithubClient(t *testing.T) {
	// Authorization Header tracking
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader = req.Header.Get("Authorization")
	}))
	defer ts.Close()

	ctx := context.Background()

	{
		// Client use `GITHUB_TOKEN`
		os.Setenv("GITHUB_TOKEN", "tokenFromEnv")

		client := getGithubClient(ctx)
		client.BaseURL, _ = url.Parse(ts.URL + "/")
		client.Repositories.GetLatestRelease(ctx, "owner", "repo")
		assert.Equal(t, "Bearer tokenFromEnv", authHeader, "token is included in request")

		os.Unsetenv("GITHUB_TOKEN")
	}

	{
		// Client doesn't use token
		client := getGithubClient(ctx)
		client.BaseURL, _ = url.Parse(ts.URL + "/")
		client.Repositories.GetLatestRelease(ctx, "owner", "repo")
		assert.Equal(t, "", authHeader, "token is not included in request")
	}
}

func TestGetGithubToken(t *testing.T) {
	{
		// Get token from environment variable
		os.Setenv("GITHUB_TOKEN", "tokenFromEnv")
		token := getGithubToken()
		assert.Equal(t, "tokenFromEnv", token)
		os.Unsetenv("GITHUB_TOKEN")
	}

	{
		// Get token from .gitconfig
		os.Setenv("GIT_CONFIG", "testdata/gitconfig")
		token := getGithubToken()
		assert.Equal(t, "tokenFromGitConfig", token)
		os.Unsetenv("GIT_CONFIG")
	}

	{
		// Cannot get token from not existing file
		os.Setenv("GIT_CONFIG", "testdata/not_exists")
		token := getGithubToken()
		assert.Equal(t, "", token)
		os.Unsetenv("GIT_CONFIG")
	}
}
