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

// Makes related environments empty,
// and returns teardown function which reset
// these variables.
func githubTestSetup() func() {
	origTokenEnv := os.Getenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_TOKEN")
	origGitConfigEnv := os.Getenv("GIT_CONFIG")
	os.Unsetenv("GIT_CONFIG")

	return func() {
		os.Setenv("GITHUB_TOKEN", origTokenEnv)
		os.Setenv("GIT_CONFIG", origGitConfigEnv)
	}
}

func TestGetGithubClient(t *testing.T) {
	teardown := githubTestSetup()
	defer teardown()

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
		os.Setenv("GIT_CONFIG", "testdata/not_exists")
		client := getGithubClient(ctx)
		client.BaseURL, _ = url.Parse(ts.URL + "/")
		client.Repositories.GetLatestRelease(ctx, "owner", "repo")
		assert.Equal(t, "", authHeader, "token is not included in request")
		os.Unsetenv("GIT_CONFIG")
	}
}

func TestGetGithubToken(t *testing.T) {
	teardown := githubTestSetup()
	defer teardown()

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
