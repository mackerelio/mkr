package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mackerelio/mkr/logger"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

// CommandPlugin is definition of mkr plugin
var CommandPlugin = cli.Command{
	Name:        "plugin",
	Usage:       "Manage mackerel plugin",
	Description: `[WIP] Manage mackerel plugin`,
	Subcommands: []cli.Command{
		{
			Name:        "install",
			Usage:       "install mackerel plugin",
			Description: `WIP`,
			Action:      doPluginInstall,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix",
					Usage: "plugin install location",
				},
			},
		},
	},
	Hidden: true,
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
	err = installByArtifact(artifactFile, filepath.Join(pluginDir, "bin"), workdir)
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

// Download plugin artifact from `url` to `workdir`,
// and returns downloaded filepath
func downloadPluginArtifact(url, workdir string) (fpath string, err error) {
	logger.Log("", fmt.Sprintf("Downloading %s", url))

	// Create request to download
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "mkr-plugin-installer/0.0.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http response not OK. code: %d, url: %s", resp.StatusCode, url)
		return "", err
	}

	// fpath is filepath where artifact will be saved
	fpath = filepath.Join(workdir, path.Base(url))

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
func installByArtifact(artifactFile, bindir, workdir string) error {
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
			return placePlugin(path, filepath.Join(bindir, name))
		}

		// `path` is a file but not plugin.
		return nil
	})
}

func looksLikePlugin(name string) bool {
	return strings.HasPrefix(name, "check-") || strings.HasPrefix(name, "mackerel-plugin-")
}

func placePlugin(src, dest string) error {
	_, err := os.Stat(dest)
	if err == nil {
		logger.Log("", fmt.Sprintf("%s already exists. Skip installing for now", dest))
		return nil
	}
	logger.Log("", fmt.Sprintf("Installing %s", dest))
	return os.Rename(src, dest)
}

type installTarget struct {
	owner        string
	repo         string
	pluginName   string
	releaseTag   string
	rawGithubURL string
}

const defaultRawGithubURL = "https://raw.githubusercontent.com"

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
		owner:        matches[1],
		repo:         matches[2],
		pluginName:   matches[3],
		releaseTag:   matches[4],
		rawGithubURL: defaultRawGithubURL,
	}
	return it, nil
}

// Make artifact's download URL
func (it *installTarget) makeDownloadURL() (string, error) {
	owner, repo, err := it.getOwnerAndRepo()
	if err != nil {
		return "", err
	}

	if it.releaseTag == "" {
		// TODO: Fetch latest release tag by github API
		return "", fmt.Errorf("not implemented")
	}

	filename := fmt.Sprintf("%s_%s_%s.zip", it.repo, runtime.GOOS, runtime.GOARCH)
	downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		owner, repo, it.releaseTag, filename)

	return downloadURL, nil
}

func (it *installTarget) getOwnerAndRepo() (string, string, error) {
	if it.owner != "" && it.repo != "" {
		return it.owner, it.repo, nil
	}

	// Get owner and repo from plugin registry
	defURL := fmt.Sprintf("%s/mackerelio/plugin-registry/master/plugins/%s.json", it.rawGithubURL, it.pluginName)
	req, err := http.NewRequest(http.MethodGet, defURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "mkr-plugin-installer/0.0.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http response not OK. code: %d, url: %s", resp.StatusCode, defURL)
		return "", "", err
	}

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

// registryDef represents one plugin definition in plugin-registry
// See Also: https://github.com/mackerelio/plugin-registry
type registryDef struct {
	Source      string `json:"source"`
	Description string `json:"description"`
}
