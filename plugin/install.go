package plugin

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mackerelio/mkr/logger"
	"github.com/mholt/archives"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
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

var commandPluginInstall = &cli.Command{
	Name:      "install",
	Usage:     "Install a plugin from github or plugin registry",
	ArgsUsage: "[--prefix <prefix>] [--overwrite] [--upgrade] <install_target>",
	Action:    doPluginInstall,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "prefix",
			Usage: fmt.Sprintf("Plugin install location. The default is %s", defaultPluginInstallLocation),
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite a plugin command in a plugin directory, even if same name command exists",
		},
		&cli.BoolFlag{
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
    - <direct_url>
          Install from specified URL.
          Supported URL schemes are http, https and file.
          Example: mkr plugin install https://github.com/mackerelio/mackerel-plugin-sample/releases/download/v0.0.3/mackerel-plugin-sample_linux_amd64.zip

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
var errSkipInstall = errors.New("skip installing for now")

// main function for mkr plugin install
func doPluginInstall(ctx context.Context, c *cli.Command) error {
	argInstallTarget := c.Args().First()
	if argInstallTarget == "" {
		return fmt.Errorf("Specify install target") // nolint
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
	workdir, err := os.MkdirTemp(filepath.Join(pluginDir, "work"), "mkr-plugin-installer-")
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while creating a work directory")
	}
	defer os.RemoveAll(workdir)

	// Download an artifact and install by it
	downloadURL, err := it.makeDownloadURL(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while making a download URL")
	}

	isMetaDataStoreEnabled := true
	meta, err := newMetaDataStore(ctx, pluginDir, it)
	if err != nil {
		if err == errDisaleMetaDataStore {
			isMetaDataStoreEnabled = false
		} else {
			return errors.Wrap(err, "Failed to prepare meta data store")
		}
	}

	overwrite := c.Bool("overwrite")
	if isMetaDataStoreEnabled && c.Bool("upgrade") {
		releaseTag, err := meta.load("release_tag")
		if err != nil {
			return errors.Wrap(err, "Failed to load release_tag")
		}
		if releaseTag == it.releaseTag {
			logger.Log("", fmt.Sprintf("release_tag %s is already installed. Skip installing for now", it.releaseTag))
			return nil
		}
		overwrite = true // force overwrite in upgrade
	}

	artifactFile, err := downloadPluginArtifact(ctx, downloadURL, workdir)
	if err != nil {
		return errors.Wrap(err, "Failed to install plugin while downloading an artifact")
	}
	err = installByArtifact(ctx, artifactFile, filepath.Join(pluginDir, "bin"), overwrite)
	if err == nil { // nolint
		if meta != nil {
			if err := meta.store("release_tag", it.releaseTag); err != nil {
				return errors.Wrap(err, "Failed to store release_tag")
			}
		}
	} else if err == errSkipInstall {
		// do not update metadata
	} else {
		return errors.Wrap(err, "Failed to install plugin while extracting and placing")
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
	err = os.MkdirAll(filepath.Join(pluginDir, "meta"), 0755)
	if err != nil {
		return "", err
	}
	return pluginDir, nil
}

// Download plugin artifact from `u`(URL) to `workdir`,
// and returns downloaded filepath
func downloadPluginArtifact(ctx context.Context, u, workdir string) (fpath string, err error) {
	logger.Log("", fmt.Sprintf("Downloading %s", u))

	// Create request to download
	resp, err := (&client{}).get(ctx, u)
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
func installByArtifact(ctx context.Context, artifactFile, bindir string, overwrite bool) error {
	fd, err := os.Open(artifactFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	format, stream, err := archives.Identify(ctx, artifactFile, fd)
	if err != nil {
		return err
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unexpected format : %s", path.Base(artifactFile))
	}

	// Look for plugin files recursively, and place those to binPath
	return ex.Extract(ctx, stream, func(ctx context.Context, info archives.FileInfo) error {
		if info.IsDir() {
			return nil
		}

		// a plugin file should be executable, and have specified name.
		name := info.Name()
		isExecutable := isWin || (info.Mode()&0111) != 0
		if isExecutable && looksLikePlugin(name) {
			fd, err := info.Open()
			if err != nil {
				return err
			}
			defer fd.Close()
			b, err := io.ReadAll(fd)
			if err != nil {
				return err
			}

			return placePlugin(b, info.Mode(), filepath.Join(bindir, name), overwrite)
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

// Create temporarily file and move it to destination
// to avoid partial write or text busy error
func safeWriteFile(dest string, data []byte, mode os.FileMode) error {
	tmpFile, err := os.CreateTemp(filepath.Dir(dest), "tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if err := tmpFile.Chmod(mode); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to set mode: %w", err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpName, dest); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func placePlugin(b []byte, mode fs.FileMode, dest string, overwrite bool) error {
	_, err := os.Stat(dest)
	if err == nil && !overwrite {
		logger.Log("", fmt.Sprintf("%s already exists. Skip installing for now", dest))
		return errSkipInstall
	}
	logger.Log("", fmt.Sprintf("Installing %s", dest))
	return safeWriteFile(dest, b, mode)
}
