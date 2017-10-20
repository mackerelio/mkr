package plugin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGithubToken(t *testing.T) {
	// Makes related environments empty,
	// and reset these after test
	origTokenEnv := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", origTokenEnv)
	os.Unsetenv("GITHUB_TOKEN")
	origGitConfigEnv := os.Getenv("GIT_CONFIG")
	defer os.Setenv("GIT_CONFIG", origGitConfigEnv)
	os.Unsetenv("GIT_CONFIG")

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
