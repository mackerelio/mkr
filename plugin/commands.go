package plugin

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

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

	_, err := parseInstallTarget(argInstallTarget)
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while setup plugin directory")
	}

	prefix, err := setupPluginDir(c.String("prefix"))
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while parsing install target")
	}

	// Create a work directory for downloading and extracting an artifact
	workdir, err := ioutil.TempDir(prefix, "work-")
	if err != nil {
		return errors.Wrap(err, "failed to install plugin while creating a work directory")
	}
	defer os.RemoveAll(workdir)

	fmt.Println("do plugin install [wip]")
	return nil
}

// Parse install target string passed from args
// example is below
// - mackerelio/mackerel-plugin-sample
// - mackerel-plugin-sample
// - mackerelio/mackerel-plugin-sample@v0.0.1
func parseInstallTarget(target string) (*installTarget, error) {
	it := &installTarget{}

	ownerRepoAndReleaseTag := strings.Split(target, "@")
	var ownerRepo string
	switch len(ownerRepoAndReleaseTag) {
	case 1:
		ownerRepo = ownerRepoAndReleaseTag[0]
	case 2:
		ownerRepo = ownerRepoAndReleaseTag[0]
		it.releaseTag = ownerRepoAndReleaseTag[1]
	default:
		return nil, fmt.Errorf("Install target is invalid: %s", target)
	}

	ownerAndRepo := strings.Split(ownerRepo, "/")
	switch len(ownerAndRepo) {
	case 1:
		it.pluginName = ownerAndRepo[0]
	case 2:
		it.owner = ownerAndRepo[0]
		it.repo = ownerAndRepo[1]
	default:
		return nil, fmt.Errorf("Install target is invalid: %s", target)
	}

	return it, nil
}

// Create a directory for plugin install
func setupPluginDir(prefix string) (string, error) {
	if prefix == "" {
		prefix = "/opt/mackerel-agent/plugins"
	}
	err := os.MkdirAll(filepath.Join(prefix, "bin"), 0755)
	if err != nil {
		return "", err
	}
	return prefix, nil
}

// Download plugin artifact from `url` to `workdir`,
// and returns downloaded filepath
func downloadPluginArtifact(url, workdir string) (fpath string, err error) {
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
	if resp.StatusCode != 200 {
		err = fmt.Errorf("http response not OK. code: %d, url: %s", resp.StatusCode, url)
		return "", err
	}

	// fpath is filepath where artifact will be saved
	fpath = filepath.Join(workdir, path.Base(url))

	// download artifact
	file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

func looksLikePlugin(name string) bool {
	return strings.HasPrefix(name, "check-") || strings.HasPrefix(name, "mackerel-plugin-")
}

type installTarget struct {
	owner      string
	repo       string
	pluginName string
	releaseTag string
}

// Make artifact's download URL
func (it *installTarget) makeDownloadURL() (string, error) {
	if it.owner != "" && it.repo != "" {
		if it.releaseTag == "" {
			// TODO: Make latest release download URL by github API
			return "", fmt.Errorf("not implemented")
		}
		filename := fmt.Sprintf("%s_%s_%s.zip", it.repo, runtime.GOOS, runtime.GOARCH)
		return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
			it.owner, it.repo, it.releaseTag, filename), nil
	}
	// TODO: Make download URL by plugin registry
	return "", fmt.Errorf("not implemented")
}
