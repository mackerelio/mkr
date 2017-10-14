package plugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func assertEqualFileContent(t *testing.T, aFile, bFile, message string) {
	aContent, err := ioutil.ReadFile(aFile)
	if err != nil {
		t.Fatal(err)
	}
	bContent, err := ioutil.ReadFile(bFile)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, aContent, bContent, message)
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
		assert.Equal(t, tc.Output, *it, "Parsing result is expected")
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
		assert.Equal(t, tc.Output, err.Error(), "error message is expected")
	}
}

func TestSetupPluginDir(t *testing.T) {
	{
		// Creating plugin dir is successful
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)

		prefix, err := setupPluginDir(tmpd)
		assert.Equal(t, tmpd, prefix, "returns default prefix directory")
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

		prefix, err := setupPluginDir(tmpd)
		assert.Equal(t, "", prefix, "returns empty string when failed")
		assert.NotNil(t, err, "error should be occured while manipulate unpermitted directory")
	}
}

func TestDownloadPluginArtifact(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()

	{
		// Response not found
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)

		fpath, err := downloadPluginArtifact(ts.URL+"/not_found.zip", tmpd)
		assert.Equal(t, "", fpath, "fpath is empty")
		assert.Contains(t, err.Error(), "http response not OK. code: 404,", "Returns correct err")
	}

	{
		// Download is finished successfully
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)

		fpath, err := downloadPluginArtifact(ts.URL+"/mackerel-plugin-sample_linux_amd64.zip", tmpd)
		assert.Equal(t, tmpd+"/mackerel-plugin-sample_linux_amd64.zip", fpath, "Returns fpath correctly")

		_, err = os.Stat(fpath)
		assert.Nil(t, err, "Downloaded file is created")

		assertEqualFileContent(t, fpath, "testdata/mackerel-plugin-sample_linux_amd64.zip", "Downloaded data is correct")
	}
}

func TestInstallByArtifact(t *testing.T) {
	{
		// Install by the artifact which has a single plugin
		bindir := tempd(t)
		defer os.RemoveAll(bindir)
		workdir := tempd(t)
		defer os.RemoveAll(workdir)

		err := installByArtifact("testdata/mackerel-plugin-sample_linux_amd64.zip", bindir, workdir)
		assert.Nil(t, err, "installByArtifact finished successfully")

		installedPath := filepath.Join(bindir, "mackerel-plugin-sample")

		fi, err := os.Stat(installedPath)
		assert.Nil(t, err, "A plugin file exists")
		assert.True(t, fi.Mode().IsRegular() && fi.Mode().Perm() == 0755, "A plugin file has execution permission")
		assertEqualFileContent(
			t,
			installedPath,
			"testdata/mackerel-plugin-sample_linux_amd64/mackerel-plugin-sample",
			"Installed plugin is valid",
		)

		// Install same name plugin, but it is skipped
		workdir2 := tempd(t)
		defer os.RemoveAll(workdir2)
		err = installByArtifact("testdata/mackerel-plugin-sample-duplicate_linux_amd64.zip", bindir, workdir2)
		assert.Nil(t, err, "installByArtifact finished successfully even if same name plugin exists")

		fi, err = os.Stat(filepath.Join(bindir, "mackerel-plugin-sample"))
		assert.Nil(t, err, "A plugin file exists")
		assertEqualFileContent(
			t,
			installedPath,
			"testdata/mackerel-plugin-sample_linux_amd64/mackerel-plugin-sample",
			"Install is skipped, so the contents is what is before",
		)
	}

	{
		// Install by the artifact which has multiple plugins
		bindir := tempd(t)
		defer os.RemoveAll(bindir)
		workdir := tempd(t)
		defer os.RemoveAll(workdir)

		installByArtifact("testdata/mackerel-plugin-sample-multi_darwin_386.zip", bindir, workdir)

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

		_, err := os.Stat(filepath.Join(bindir, "mackerel-plugin-not-executable"))
		assert.NotNil(t, err, "mackerel-plugin-not-executable is not installed")
		_, err = os.Stat(filepath.Join(bindir, "not-mackerel-plugin-sample"))
		assert.NotNil(t, err, "not-mackerel-plugin-sample is not installed")
	}
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
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.LooksLikePlugin, looksLikePlugin(tc.Name))
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
			fmt.Sprintf("https://github.com/mackerelio/mackerel-plugin-sample/releases/download/v0.1.0/mackerel-plugin-sample_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
			url,
			"Download URL is made correctly",
		)
	}
}
