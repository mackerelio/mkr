package plugin

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	gitconfig "github.com/tcnksm/go-gitconfig"
	"golang.org/x/oauth2"
)

// Get github client having github token.
func getGithubClient(ctx context.Context) *github.Client {
	var oauthClient *http.Client
	if token := getGithubToken(); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: getGithubToken()},
		)
		oauthClient = oauth2.NewClient(ctx, ts)
	}
	return github.NewClient(oauthClient)
}

// Get github token from environment variables, or github.token in gitconfig file
func getGithubToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}
	token, _ = gitconfig.GithubToken()
	return token
}
