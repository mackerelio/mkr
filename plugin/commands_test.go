package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tempd(t *testing.T) string {
	tmpd, err := ioutil.TempDir("", "mkr-plugin-install")
	if err != nil {
		t.Fatal(err)
	}
	return tmpd
}

func TestSetupPluginDir(t *testing.T) {
	{
		// Creating plugin dir is successful
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)
		err := setupPluginDir(tmpd)

		assert.Nil(t, err, "setup finished successfully")
		fi, err := os.Stat(filepath.Join(tmpd, "bin"))
		if assert.Nil(t, err) {
			assert.True(t, fi.IsDir(), "plugin directory is created")
		}
	}

	{
		// Creating plugin dir is failed because of directory's permission
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)
		err := os.Chmod(tmpd, 0500)
		assert.Nil(t, err, "chmod finished successfully")

		err = setupPluginDir(tmpd)
		assert.NotNil(t, err, "error should be occured while manipulate unpermitted directory")
	}
}

func TestParseInstallTarget(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output installTarget
	}{
		{
			Name:  "Plugin name only",
			Input: "mackerel-plugin-sample",
			Output: installTarget{
				pluginName: "mackerel-plugin-sample",
			},
		},
		{
			Name:  "Plugin name and release tag",
			Input: "mackerel-plugin-sample@v0.0.1",
			Output: installTarget{
				pluginName: "mackerel-plugin-sample",
				releaseTag: "v0.0.1",
			},
		},
		{
			Name:  "Owner and repo",
			Input: "mackerelio/mackerel-plugin-sample",
			Output: installTarget{
				owner: "mackerelio",
				repo:  "mackerel-plugin-sample",
			},
		},
		{
			Name:  "Owner and repo with release tag",
			Input: "mackerelio/mackerel-plugin-sample@v1.0.1",
			Output: installTarget{
				owner:      "mackerelio",
				repo:       "mackerel-plugin-sample",
				releaseTag: "v1.0.1",
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("testing: %s\n", tc.Name)
		it, err := parseInstallTarget(tc.Input)
		assert.Nil(t, err, "error does not occur while parseInstallTarget")
		assert.Equal(t, *it, tc.Output, "Parsing result is expected")
	}
}

func TestParseInstallTarget_error(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output string
	}{
		{
			Name:   "Too many @",
			Input:  "mackerel-plugin-sample@v0.0.1@v0.1.0",
			Output: "Install target is invalid: mackerel-plugin-sample@v0.0.1@v0.1.0",
		},
		{
			Name:   "Too many /",
			Input:  "mackerelio/hatena/mackerel-plugin-sample",
			Output: "Install target is invalid: mackerelio/hatena/mackerel-plugin-sample",
		},
	}

	for _, tc := range testCases {
		t.Logf("testing: %s\n", tc.Name)
		_, err := parseInstallTarget(tc.Input)
		assert.NotNil(t, err, "parseInstallTarget returns err when invalid target string is passed")
		assert.Equal(t, err.Error(), tc.Output, "error message is expected")
	}
}

func TestInstallTargetMakeDownloadURL(t *testing.T) {
	{
		t.Logf("Make download URL with owner, repo and releaseTag")
		it := &installTarget{
			owner:      "mackerelio",
			repo:       "mackerel-plugin-sample",
			releaseTag: "v0.1.0",
		}
		url, err := it.makeDownloadURL()
		assert.Nil(t, err, "makeDownloadURL is successful")
		assert.Equal(
			t,
			url,
			fmt.Sprintf("https://github.com/mackerelio/mackerel-plugin-sample/releases/download/v0.1.0/mackerel-plugin-sample_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
			"Download URL is made correctly",
		)
	}
}
