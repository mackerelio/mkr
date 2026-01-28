package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type installTarget struct {
	owner      string
	repo       string
	pluginName string
	releaseTag string
	directURL  string

	// fields for testing
	rawGithubURL string
	apiGithubURL string
	githubURL    string
}

const (
	defaultRawGithubURL = "https://raw.githubusercontent.com"
	defaultAPIGithubURL = "https://api.github.com"
	defaultGithubURL    = "https://github.com"
)

// the pattern of installTarget string
var (
	// (?:<plugin_name>|<owner>/<repo>)(?:@<releaseTag>)?
	targetReg = regexp.MustCompile(`^(?:([^@/]+)/([^@/]+)|([^@/]+))(?:@(.+))?$`)
	urlReg    = regexp.MustCompile(`^(?:https?|file)://`)
)

// Parse install target string, and construct installTarget
// example is below
// - mackerelio/mackerel-plugin-sample
// - mackerel-plugin-sample
// - mackerelio/mackerel-plugin-sample@v0.0.1
// - https://mackerel.io/mackerel-plugin-sample_linux_amd64.zip
// - file:///path/to/mackerel-plugin-sample_linux_amd64.zip
func newInstallTargetFromString(target string) (*installTarget, error) {
	if urlReg.MatchString(target) {
		return &installTarget{
			directURL: target,
		}, nil
	}

	matches := targetReg.FindStringSubmatch(target)
	if len(matches) != 5 {
		return nil, fmt.Errorf("Install target is invalid: %s", target) // nolint
	}

	it := &installTarget{
		owner:      matches[1],
		repo:       matches[2],
		pluginName: matches[3],
		releaseTag: matches[4],
	}
	return it, nil
}

// Make artifact's download URL
func (it *installTarget) makeDownloadURL(ctx context.Context) (string, error) {
	if it.directURL != "" {
		return it.directURL, nil
	}

	owner, repo, err := it.getOwnerAndRepo(ctx)
	if err != nil {
		return "", err
	}

	releaseTag, err := it.getTagFromReleasesURL(ctx, owner, repo)
	if releaseTag == "" || err != nil {
		releaseTag, err = it.getReleaseTag(owner, repo)
		if err != nil {
			return "", err
		}
	}

	filename := fmt.Sprintf("%s_%s_%s.zip", url.PathEscape(repo), runtime.GOOS, runtime.GOARCH)
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		url.PathEscape(owner),
		url.PathEscape(repo),
		url.PathEscape(releaseTag),
		filename,
	)

	return downloadURL, nil
}

func (it *installTarget) getOwnerAndRepo(ctx context.Context) (string, string, error) {
	if it.owner != "" && it.repo != "" {
		return it.owner, it.repo, nil
	}

	// if directURL is specified, target doesn't have owner and repo
	if it.directURL != "" {
		return "", "", fmt.Errorf("owner and repo are not found because directURL is specified")
	}

	// Get owner and repo from plugin registry
	defURL := fmt.Sprintf(
		"%s/mackerelio/plugin-registry/master/plugins/%s.json",
		it.getRawGithubURL(),
		url.PathEscape(it.pluginName),
	)
	resp, err := (&client{}).get(ctx, defURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var def registryDef
	err = json.NewDecoder(resp.Body).Decode(&def)
	if err != nil {
		return "", "", err
	}

	ownerAndRepo := strings.Split(def.Source, "/")
	if len(ownerAndRepo) != 2 {
		return "", "", fmt.Errorf("source definition is invalid")
	}

	// Cache owner and repo
	it.owner = ownerAndRepo[0]
	it.repo = ownerAndRepo[1]

	return it.owner, it.repo, nil
}

func (it *installTarget) getReleaseTag(owner, repo string) (string, error) {
	if it.releaseTag != "" {
		return it.releaseTag, nil
	}

	// Get latest release tag from Github API
	ctx := context.Background()
	client := getGithubClient(ctx)
	client.BaseURL = it.getAPIGithubURL()

	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	// Cache releaseTag
	it.releaseTag = release.GetTagName()
	return it.releaseTag, nil
}

func (it *installTarget) getTagFromReleasesURL(ctx context.Context, owner, repo string) (string, error) {
	if it.releaseTag != "" {
		return it.releaseTag, nil
	}

	latestURL := fmt.Sprintf("%s/%s/%s/releases/latest", it.getGithubURL(), owner, repo)

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, latestURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	defer resp.Body.Close()

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("no Location header found (status=%d)", resp.StatusCode)
	}

	u, err := resp.Request.URL.Parse(loc)
	if err != nil {
		return "", err
	}

	tag, err := extractTagFromReleasesURL(u)
	if err != nil {
		return "", err
	}
	it.releaseTag = tag
	return it.releaseTag, nil
}

func extractTagFromReleasesURL(u *url.URL) (string, error) {
	// /<owner>/<repo>/releases/tag/v1.2.3
	segments := strings.Split(strings.Trim(u.EscapedPath(), "/"), "/")
	for i := 0; i < len(segments)-1; i++ {
		if segments[i] == "releases" && i+1 < len(segments) && segments[i+1] == "tag" {
			// Ends with a tag
			raw := path.Base(u.EscapedPath())
			tag, err := url.PathUnescape(raw)
			if err != nil {
				return "", err
			}
			return tag, nil
		}
	}
	return "", fmt.Errorf("not a releases/tag URL")
}

func (it *installTarget) getRawGithubURL() string {
	if it.rawGithubURL != "" {
		return it.rawGithubURL
	}
	return defaultRawGithubURL
}

func (it *installTarget) getGithubURL() string {
	if it.githubURL != "" {
		return it.githubURL
	}
	return defaultGithubURL
}

// Returns URL object which Github Client.BaseURL can receive as it is
func (it *installTarget) getAPIGithubURL() *url.URL {
	u := defaultAPIGithubURL
	if it.apiGithubURL != "" {
		u = it.apiGithubURL
	}
	// Ignore err because apiGithubURL is specified only internally
	apiURL, _ := url.Parse(u + "/") // trailing `/` is required for BaseURL
	return apiURL
}

// registryDef represents one plugin definition in plugin-registry
// See Also: https://github.com/mackerelio/plugin-registry
type registryDef struct {
	Source      string `json:"source"`
	Description string `json:"description"`
}
