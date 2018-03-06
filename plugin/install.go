package plugin

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mackerelio/mkr/logger"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

// The reason why an immediate function, not `init()` is used here is that
// the `defaultPluginInstallLocation` is used in following `commandPluginInstall`
// assignment. Top level variable assignment is executed before `init()`.
var defaultPluginInstallLocation = func() string {
	if runtime.GOOS != "windows" {
		return "/opt/mackerel-agent/plugins"
	}
	path, err := os.Executable()
	logger.DieIf(err)
	return filepath.Join(filepath.Dir(path), "plugins")
}()

var commandPluginInstall = cli.Command{
	Name:      "install",
	Usage:     "Install a plugin from github or plugin registry",
	ArgsUsage: "[--prefix <prefix>] [--overwrite] <install_target>",
	Action:    doPluginInstall,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "prefix",
			Usage: fmt.Sprintf("Plugin install location. The default is %s", defaultPluginInstallLocation),
		},
		cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite a plugin command in a plugin directory, even if same name command exists",
		},
		cli.BoolFlag{
			Name:  "upgrade",
			Usage: "Upgrade a plugin command in a plugin directory only when a release_tag is modified",
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

    Please refer to following documents for detail.
    - Using mkr plugin install
      https://mackerel.io/docs/entry/advanced/install-plugin-by-mkr
    - Creating plugins supported with mkr plugin install
      https://mackerel.io/docs/entry/advanced/make-plugin-corresponding-to-installer
`,
}

var isWin = runtime.GOOS == "windows"

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

	if c.Bool("upgrade") {
		install, err := shouldInstall(pluginDir, it)
		if err != nil {
			return errors.Wrap(err, "Failed to detect plugin should be installed")
		}
		if !install {
			logger.Log("", fmt.Sprintf("release_tag %s already exists. Skip installing for now", it.releaseTag))
			return nil
		}
	}

	artifactFile, err := downloadPluginArtifact(downloadURL, workdir)
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while downloading an artifact")
	}
	err = installByArtifact(artifactFile, filepath.Join(pluginDir, "bin"), workdir, c.Bool("overwrite"))
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while extracting and placing")
	}

	err = storeReleaseTag(pluginDir, it)
	if err != nil {
		return errors.Wrap(err, "Failed to store plugin release tag")
	}

	logger.Log("", fmt.Sprintf("Successfully installed %s", argInstallTarget))
	return nil
}

// Create a directory for plugin install
func setupPluginDir(pluginDir string) (string, error) {
	if pluginDir == "" {
		pluginDir = defaultPluginInstallLocation
	}
	err := os.MkdirAll(filepath.Join(pluginDir, "bin"), 0755)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(filepath.Join(pluginDir, "work"), 0755)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(filepath.Join(pluginDir, "tags"), 0755)
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
	fn := archiver.Zip.Open
	if strings.HasSuffix(artifactFile, ".tar.gz") || strings.HasSuffix(artifactFile, ".tgz") {
		fn = archiver.TarGz.Open
	}
	err := fn(artifactFile, workdir)
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
		isExecutable := isWin || (info.Mode()&0111) != 0
		if isExecutable && looksLikePlugin(name) {
			return placePlugin(path, filepath.Join(bindir, name), overwrite)
		}
		// `path` is a file but not plugin.
		return nil
	})
}

func looksLikePlugin(name string) bool {
	if strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz") {
		return false
	}
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

func shouldInstall(pluginDir string, it *installTarget) (bool, error) {
	dir := filepath.Join(pluginDir, "tags", it.owner)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, err
	}
	filename := filepath.Join(dir, it.repo)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if string(b) == it.releaseTag {
		return false, nil
	}

	return true, nil
}

func storeReleaseTag(pluginDir string, it *installTarget) error {
	dir := filepath.Join(pluginDir, "tags", it.owner)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	filename := filepath.Join(dir, it.repo)
	return ioutil.WriteFile(filename, []byte(it.releaseTag), 0644)
}
