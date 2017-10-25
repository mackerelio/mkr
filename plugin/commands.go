package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"net/url"

	"github.com/mackerelio/mkr/logger"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

// CommandPlugin is definition of mkr plugin
var CommandPlugin = cli.Command{
	Name:  "plugin",
	Usage: "Manage mackerel plugin",
	Description: `
    Manage mackerel plugin.  For example, you can install a mackerel plugin and
    check plugin by "mkr plugin install".
`,
	Subcommands: []cli.Command{
		{
			Name:      "install",
			Usage:     "Install a plugin from github or plugin registry",
			ArgsUsage: "[--prefix <prefix>] [--overwrite] <install_target>",
			Action:    doPluginInstall,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix",
					Usage: "Plugin install location. The default is /opt/mackerel-agent/plugins",
				},
				cli.BoolFlag{
					Name:  "overwrite",
					Usage: "Overwrite a plugin command in a plugin directory, even if same name command exists",
				},
			},
			Description: `
    Install a mackerel plugin and a check plugin from github or plugin registry.
    To install by mkr, a plugin has to be released to Github Releases in specification format.

    <install_target> is:
    - <owner>/<repo>[@<release_tag>]
          Install from specified github owner, repository, and Github Releases tag.
          If you omit <release_tag>, the installer install from latest Github Release.
          Example: mkr plugin install mackerelio/mackerel-plugin-sample@v0.0.1
    - <plugin_name>[@<release_tag]
          Install from plugin registry.
          You can find available plugins in https://github.com/mackerelio/plugin-registry
          Example: mkr plugin install mackerel-plugin-sample

    The installer uses Github API to find the latest release.  Please set a github token to
    GITHUB_TOKEN environment variable, or to github.token in .gitconfig.
    Otherwise, installation sometimes fails because of Github API Rate Limit.

    If you want to use the plugin installer by a server provisioning tool,
    we recommend you to specify <release_tag> explicitly.
    If you specify <release_tag>, the installer doesn't use Github API,
    so Github API Rate Limit error doesn't occur.
`,
		},
	},
}

// main function for mkr plugin install
func doPluginInstall(c *cli.Context) error {
	argInstallTarget := c.Args().First()
	if argInstallTarget == "" {
		return fmt.Errorf("Specify install target")
	}

	it, err := newInstallTargetFromString(argInstallTarget)
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while parsing install target")
	}

	pluginDir, err := setupPluginDir(c.String("prefix"))
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while setup plugin directory")
	}

	// Create a work directory for downloading and extracting an artifact
	workdir, err := ioutil.TempDir(filepath.Join(pluginDir, "work"), "mkr-plugin-installer-")
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while creating a work directory")
	}
	defer os.RemoveAll(workdir)

	// Download an artifact and install by it
	downloadURL, err := it.makeDownloadURL()
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while making a download URL")
	}
	artifactFile, err := downloadPluginArtifact(downloadURL, workdir)
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while downloading an artifact")
	}
	err = installByArtifact(artifactFile, filepath.Join(pluginDir, "bin"), workdir, c.Bool("overwrite"))
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while extracting and placing")
	}

	logger.Log("", fmt.Sprintf("Successfully installed %s", argInstallTarget))
	return nil
}

// Create a directory for plugin install
func setupPluginDir(pluginDir string) (string, error) {
	if pluginDir == "" {
		pluginDir = "/opt/mackerel-agent/plugins"
	}
	err := os.MkdirAll(filepath.Join(pluginDir, "bin"), 0755)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(filepath.Join(pluginDir, "work"), 0755)
	if err != nil {
		return "", err
	}
	return pluginDir, nil
}

// Download plugin artifact from `u`(URL) to `workdir`,
// and returns downloaded filepath
func downloadPluginArtifact(u, workdir string) (fpath string, err error) {
	logger.Log("", fmt.Sprintf("Downloading %s", u))

	// Create request to download
	resp, err := (&client{}).get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// fpath is filepath where artifact will be saved
	fpath = filepath.Join(workdir, path.Base(u))

	// download artifact
	file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return fpath, nil
}

// Extract artifact and install plugin
func installByArtifact(artifactFile, bindir, workdir string, overwrite bool) error {
	// unzip artifact to work directory
	err := archiver.Zip.Open(artifactFile, workdir)
	if err != nil {
		return err
	}

	// Look for plugin files recursively, and place those to binPath
	return filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// a plugin file should be executable, and have specified name.
		name := info.Name()
		isExecutable := (info.Mode() & 0111) != 0
		if isExecutable && looksLikePlugin(name) {
			return placePlugin(path, filepath.Join(bindir, name), overwrite)
		}

		// `path` is a file but not plugin.
		return nil
	})
}

func looksLikePlugin(name string) bool {
	return strings.HasPrefix(name, "check-") || strings.HasPrefix(name, "mackerel-plugin-")
}

func placePlugin(src, dest string, overwrite bool) error {
	_, err := os.Stat(dest)
	if err == nil && !overwrite {
		logger.Log("", fmt.Sprintf("%s already exists. Skip installing for now", dest))
		return nil
	}
	logger.Log("", fmt.Sprintf("Installing %s", dest))
	return os.Rename(src, dest)
}

type installTarget struct {
	owner      string
	repo       string
	pluginName string
	releaseTag string

	// fields for testing
	rawGithubURL string
	apiGithubURL string
}

const (
	defaultRawGithubURL = "https://raw.githubusercontent.com"
	defaultAPIGithubURL = "https://api.github.com"
)

// the pattern of installTarget string
// (?:<plugin_name>|<owner>/<repo>)(?:@<releaseTag>)?
var targetReg = regexp.MustCompile(`^(?:([^@/]+)/([^@/]+)|([^@/]+))(?:@(.+))?$`)

// Parse install target string, and construct installTarget
// example is below
// - mackerelio/mackerel-plugin-sample
// - mackerel-plugin-sample
// - mackerelio/mackerel-plugin-sample@v0.0.1
func newInstallTargetFromString(target string) (*installTarget, error) {
	matches := targetReg.FindStringSubmatch(target)
	if len(matches) != 5 {
		return nil, fmt.Errorf("Install target is invalid: %s", target)
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
func (it *installTarget) makeDownloadURL() (string, error) {
	owner, repo, err := it.getOwnerAndRepo()
	if err != nil {
		return "", err
	}

	releaseTag, err := it.getReleaseTag(owner, repo)
	if err != nil {
		return "", err
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

func (it *installTarget) getOwnerAndRepo() (string, string, error) {
	if it.owner != "" && it.repo != "" {
		return it.owner, it.repo, nil
	}

	// Get owner and repo from plugin registry
	defURL := fmt.Sprintf(
		"%s/mackerelio/plugin-registry/master/plugins/%s.json",
		it.getRawGithubURL(),
		url.PathEscape(it.pluginName),
	)
	resp, err := (&client{}).get(defURL)
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

func (it *installTarget) getRawGithubURL() string {
	if it.rawGithubURL != "" {
		return it.rawGithubURL
	}
	return defaultRawGithubURL
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
