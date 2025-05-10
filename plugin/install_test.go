package plugin

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func assertEqualFileContent(t *testing.T, aFile, bFile, message string) {
	aContent, err := os.ReadFile(aFile)
	if err != nil {
		t.Fatal(err)
	}
	bContent, err := os.ReadFile(bFile)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, aContent, bContent, message)
}

func TestSetupPluginDir(t *testing.T) {
	t.Run("Creating plugin dir is successful", func(t *testing.T) {
		tmpd := t.TempDir()

		pluginDir, err := setupPluginDir(tmpd)
		assert.Equal(t, tmpd, pluginDir, "returns default plugin directory")
		assert.Nil(t, err, "setup finished successfully")

		fi, err := os.Stat(filepath.Join(tmpd, "bin"))
		if assert.Nil(t, err) {
			assert.True(t, fi.IsDir(), "plugin bin directory is created")
		}

		fi, err = os.Stat(filepath.Join(tmpd, "work"))
		if assert.Nil(t, err) {
			assert.True(t, fi.IsDir(), "plugin work directory is created")
		}

		fi, err = os.Stat(filepath.Join(tmpd, "meta"))
		if assert.Nil(t, err) {
			assert.True(t, fi.IsDir(), "plugin meta directory is created")
		}
	})

	t.Run("Creating plugin dir is failed because of directory's permission", func(t *testing.T) {
		if isWin {
			t.Skip("skipping test on windows")
		}
		tmpd := t.TempDir()
		err := os.Chmod(tmpd, 0500)
		assert.Nil(t, err, "chmod finished successfully")

		pluginDir, err := setupPluginDir(tmpd)
		assert.Equal(t, "", pluginDir, "returns empty string when failed")
		assert.NotNil(t, err, "error should be occured while manipulate unpermitted directory")
	})
}

func TestDownloadPluginArtifact(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()

	t.Run("Response not found", func(t *testing.T) {
		tmpd := t.TempDir()

		fpath, err := downloadPluginArtifact(ts.URL+"/not_found.zip", tmpd)
		assert.Equal(t, "", fpath, "fpath is empty")
		assert.Contains(t, err.Error(), "http response not OK. code: 404,", "Returns correct err")
	})

	t.Run("Download is finished successfully", func(t *testing.T) {
		tmpd := t.TempDir()

		fpath, err := downloadPluginArtifact(ts.URL+"/mackerel-plugin-sample_linux_amd64.zip", tmpd)
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(tmpd, "/mackerel-plugin-sample_linux_amd64.zip"), fpath, "Returns fpath correctly")

		_, err = os.Stat(fpath)
		assert.Nil(t, err, "Downloaded file is created")

		assertEqualFileContent(t, fpath, "testdata/mackerel-plugin-sample_linux_amd64.zip", "Downloaded data is correct")
	})
}

func TestInstallByArtifact(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}
	{
		bindir := t.TempDir()
		installedPath := filepath.Join(bindir, "mackerel-plugin-sample")

		t.Run("Install by the artifact which has a single plugin", func(t *testing.T) {
			err := installByArtifact("testdata/mackerel-plugin-sample_linux_amd64.zip", bindir, false)
			assert.Nil(t, err, "installByArtifact finished successfully")

			fi, err := os.Stat(installedPath)
			assert.Nil(t, err, "A plugin file exists")
			if !isWin {
				assert.True(t, fi.Mode().IsRegular() && fi.Mode().Perm() == 0755, "A plugin file has execution permission")
			}
			assertEqualFileContent(
				t,
				installedPath,
				"testdata/mackerel-plugin-sample_linux_amd64/mackerel-plugin-sample",
				"Installed plugin is valid",
			)
		})

		t.Run("Install same name plugin, but it is skipped", func(t *testing.T) {
			err := installByArtifact("testdata/mackerel-plugin-sample-duplicate_linux_amd64.zip", bindir, false)
			assert.ErrorIs(t, err, errSkipInstall, "installByArtifact finished successfully even if same name plugin exists")

			_, err = os.Stat(installedPath)
			assert.Nil(t, err, "A plugin file exists")
			assertEqualFileContent(
				t,
				installedPath,
				"testdata/mackerel-plugin-sample_linux_amd64/mackerel-plugin-sample",
				"Install is skipped, so the contents is what is before",
			)
		})

		t.Run("Install same name plugin with overwrite option", func(t *testing.T) {
			err := installByArtifact("testdata/mackerel-plugin-sample-duplicate_linux_amd64.zip", bindir, true)
			assert.Nil(t, err, "installByArtifact finished successfully")
			assertEqualFileContent(
				t,
				installedPath,
				"testdata/mackerel-plugin-sample-duplicate_linux_amd64/mackerel-plugin-sample",
				"a plugin is installed with overwrite option, so the contents is overwritten",
			)
		})
	}

	t.Run("tgz", func(*testing.T) {
		bindir := t.TempDir()

		err := installByArtifact("testdata/mackerel-plugin-sample_linux_amd64.tar.gz", bindir, false)
		assert.Nil(t, err, "installByArtifact finished successfully")

		installedPath := filepath.Join(bindir, "mackerel-plugin-sample")

		fi, err := os.Stat(installedPath)
		assert.Nil(t, err, "A plugin file exists")
		if !isWin {
			assert.True(t, fi.Mode().IsRegular() && fi.Mode().Perm() == 0755, "A plugin file has execution permission")
		}
		assertEqualFileContent(
			t,
			installedPath,
			"testdata/mackerel-plugin-sample_linux_amd64/mackerel-plugin-sample",
			"Installed plugin is valid",
		)
	})

	t.Run("Install by the artifact which has multiple plugins", func(t *testing.T) {
		bindir := t.TempDir()

		err := installByArtifact("testdata/mackerel-plugin-sample-multi_darwin_386.zip", bindir, false)
		assert.Nil(t, err, "installByArtifact finished successfully")

		// check-sample, mackerel-plugin-sample-multi-1 and plugins/mackerel-plugin-sample-multi-2
		// are installed.  But followings are not installed
		// - mackerel-plugin-non-executable: does not have execution permission
		// - not-mackerel-plugin-sample: does not has plugin file name
		assertEqualFileContent(t,
			filepath.Join(bindir, "check-sample"),
			"testdata/mackerel-plugin-sample-multi_darwin_386/check-sample",
			"check-sample is installed",
		)
		assertEqualFileContent(t,
			filepath.Join(bindir, "mackerel-plugin-sample-multi-1"),
			"testdata/mackerel-plugin-sample-multi_darwin_386/mackerel-plugin-sample-multi-1",
			"mackerel-plugin-sample-multi-1 is installed",
		)
		assertEqualFileContent(t,
			filepath.Join(bindir, "mackerel-plugin-sample-multi-2"),
			"testdata/mackerel-plugin-sample-multi_darwin_386/plugins/mackerel-plugin-sample-multi-2",
			"mackerel-plugin-sample-multi-2 is installed",
		)

		_, err = os.Stat(filepath.Join(bindir, "mackerel-plugin-not-executable"))
		assert.NotNil(t, err, "mackerel-plugin-not-executable is not installed")
		_, err = os.Stat(filepath.Join(bindir, "not-mackerel-plugin-sample"))
		assert.NotNil(t, err, "not-mackerel-plugin-sample is not installed")
	})
}

func newPluginInstallContext(t testing.TB, target, prefix string, overwrite bool) *cli.Context {
	t.Helper()
	fs := flag.NewFlagSet("name", flag.ContinueOnError)
	for _, f := range commandPluginInstall.Flags {
		f.Apply(fs)
	}
	argv := []string{}
	if prefix != "" {
		argv = append(argv, fmt.Sprintf("-prefix=%s", prefix))
	}
	if overwrite {
		argv = append(argv, "-overwrite")
	}
	if target != "" {
		argv = append(argv, target)
	}
	if err := fs.Parse(argv); err != nil {
		t.Fatal(err)
	}
	return cli.NewContext(nil, fs, nil)
}

func TestDoPluginInstall(t *testing.T) {
	t.Run("specify URL directly", func(t *testing.T) {
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
		defer ts.Close()
		tmpd := t.TempDir()

		ctx := newPluginInstallContext(t, ts.URL+"/mackerel-plugin-sample_linux_amd64.zip", tmpd, false)
		err := doPluginInstall(ctx)
		assert.Nil(t, err, "sample plugin is succesfully installed")

		fpath := filepath.Join(tmpd, "bin", "mackerel-plugin-sample")
		_, err = os.Stat(fpath)
		assert.Nil(t, err, "sample plugin is successfully installed and located")
	})

	t.Run("file: scheme URL", func(t *testing.T) {
		if isWin {
			t.Skip("skipping on windows")
		}
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		fpath := filepath.Join(cwd, "testdata", "mackerel-plugin-sample_linux_amd64.zip")
		fpath = filepath.ToSlash(fpath) // care windows
		scheme := "file://"
		if !strings.HasPrefix(fpath, "/") {
			// care windows drive letter
			scheme += "/"
		}

		tmpd := t.TempDir()

		ctx := newPluginInstallContext(t, scheme+fpath, tmpd, false)
		err = doPluginInstall(ctx)
		assert.Nil(t, err, "sample plugin is succesfully installed")

		plugPath := filepath.Join(tmpd, "bin", "mackerel-plugin-sample")
		_, err = os.Stat(plugPath)
		assert.Nil(t, err, "sample plugin is successfully installed and located")
	})
}

func TestLooksLikePlugin(t *testing.T) {
	testCases := []struct {
		Name            string
		LooksLikePlugin bool
	}{
		{"mackerel-plugin-sample", true},
		{"mackerel-plugin-hoge_sample1", true},
		{"check-sample", true},
		{"check-hoge-sample", true},
		{"mackerel-sample", false},
		{"hoge-mackerel-plugin-sample", false},
		{"hoge-check-sample", false},
		{"wrong-sample", false},
		{"mackerel-plugin-sample.zip", false},
		{"mackerel-plugin-sample.tgz", false},
		{"check-sample.tar.gz", false},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.LooksLikePlugin, looksLikePlugin(tc.Name))
	}
}
