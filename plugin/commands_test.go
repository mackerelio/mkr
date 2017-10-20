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

	"encoding/json"

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

func TestSetupPluginDir(t *testing.T) {
	{
		// Creating plugin dir is successful
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)

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
	}

	{
		// Creating plugin dir is failed because of directory's permission
		tmpd := tempd(t)
		defer os.RemoveAll(tmpd)
		err := os.Chmod(tmpd, 0500)
		assert.Nil(t, err, "chmod finished successfully")

		pluginDir, err := setupPluginDir(tmpd)
		assert.Equal(t, "", pluginDir, "returns empty string when failed")
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

func TestNewInstallTargetFromString(t *testing.T) {
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
		{
			Name:  "Owner and repo with release tag(which has / and @)",
			Input: "mackerelio/mackerel-plugin-sample@v1.0.1/hoge@fuga",
			Output: installTarget{
				owner:      "mackerelio",
				repo:       "mackerel-plugin-sample",
				releaseTag: "v1.0.1/hoge@fuga",
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("testing: %s\n", tc.Name)
		it, err := newInstallTargetFromString(tc.Input)
		assert.Nil(t, err, "error does not occur while newInstallTargetFromString")
		assert.Equal(t, tc.Output, *it, "Parsing result is expected")
	}
}

func TestNewInstallTargetFromString_error(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output string
	}{
		{
			Name:   "Empty String",
			Input:  "",
			Output: "Install target is invalid: ",
		},
		{
			Name:   "Too many /",
			Input:  "mackerelio/hatena/mackerel-plugin-sample",
			Output: "Install target is invalid: mackerelio/hatena/mackerel-plugin-sample",
		},
		{
			Name:   "End with /",
			Input:  "mackerelio/",
			Output: "Install target is invalid: mackerelio/",
		},
		{
			Name:   "Start with /",
			Input:  "/mackerel-plugin-sample",
			Output: "Install target is invalid: /mackerel-plugin-sample",
		},
		{
			Name:   "Only release tag",
			Input:  "@v0.0.1",
			Output: "Install target is invalid: @v0.0.1",
		},
		{
			Name:   "End with @",
			Input:  "hoge/fuga@",
			Output: "Install target is invalid: hoge/fuga@",
		},
	}

	for _, tc := range testCases {
		t.Logf("testing: %s\n", tc.Name)
		_, err := newInstallTargetFromString(tc.Input)
		assert.NotNil(t, err, "newInstallTargetFromString returns err when invalid target string is passed")
		assert.Equal(t, tc.Output, err.Error(), "error message is expected")
	}
}

func TestInstallTargetMakeDownloadURL(t *testing.T) {
	{
		// Make download URL for `<owner>/<repo>@<releaseTag>`
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

	{
		// Make download URL for `<pluginName>@<releaseTag>`
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/mackerelio/plugin-registry/master/plugins/mackerel-plugin-hoge.json",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"description": "hoge mackerel plugin", "source": "owner-1/mackerel-plugin-hoge"}`)
			},
		)
		rawGithubServer := httptest.NewServer(mux)
		defer rawGithubServer.Close()

		it := &installTarget{
			pluginName:   "mackerel-plugin-hoge",
			releaseTag:   "v1.2.3",
			rawGithubURL: rawGithubServer.URL,
		}
		url, err := it.makeDownloadURL()
		assert.NoError(t, err, "makeDownloadURL is successful")
		assert.Equal(
			t,
			fmt.Sprintf("https://github.com/owner-1/mackerel-plugin-hoge/releases/download/v1.2.3/mackerel-plugin-hoge_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
			url,
			"Download URL is made correctly",
		)

		// Make download URL with pluginName which is not defined in registry
		it = &installTarget{
			pluginName:   "mackerel-plugin-fuga",
			releaseTag:   "v1.2.3",
			rawGithubURL: rawGithubServer.URL,
		}
		_, err = it.makeDownloadURL()
		assert.Error(t, err, "makeDownloadURL is failed")
	}

	{
		// Make download URL for `<owner>/<repo>` (latest release)
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/owner1/check-repo1/releases/latest", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"tag_name": "1.01"}`)
		})
		githubAPIServer := httptest.NewServer(mux)
		defer githubAPIServer.Close()

		it := &installTarget{
			owner:        "owner1",
			repo:         "check-repo1",
			githubAPIURL: githubAPIServer.URL,
		}
		url, err := it.makeDownloadURL()
		assert.NoError(t, err, "makeDownloadURL is successful")
		assert.Equal(
			t,
			fmt.Sprintf("https://github.com/owner1/check-repo1/releases/download/1.01/check-repo1_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
			url,
			"Download URL is made correctly",
		)

		// Latest release is not found
		it = &installTarget{
			owner:        "owner1",
			repo:         "check-not-found",
			githubAPIURL: githubAPIServer.URL,
		}
		_, err = it.makeDownloadURL()
		assert.Error(t, err, "makeDownloadURL is failed")
	}

	{
		// Make download URL for `<pluginName>`
		muxAPI := http.NewServeMux()
		muxAPI.HandleFunc("/repos/owner1/mackerel-plugin-repo1/releases/latest", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"tag_name": "release/v0.5.1"}`)
		})
		githubAPIServer := httptest.NewServer(muxAPI)
		defer githubAPIServer.Close()

		muxRaw := http.NewServeMux()
		muxRaw.HandleFunc(
			"/mackerelio/plugin-registry/master/plugins/mackerel-plugin-repo1.json",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"description": "mackerel plugin", "source": "owner1/mackerel-plugin-repo1"}`)
			},
		)
		rawGithubServer := httptest.NewServer(muxRaw)
		defer rawGithubServer.Close()

		it := &installTarget{
			pluginName:   "mackerel-plugin-repo1",
			githubAPIURL: githubAPIServer.URL,
			rawGithubURL: rawGithubServer.URL,
		}

		url, err := it.makeDownloadURL()
		assert.NoError(t, err, "makeDownloadURL is successful")
		assert.Equal(
			t,
			fmt.Sprintf("https://github.com/owner1/mackerel-plugin-repo1/releases/download/release%%2Fv0.5.1/mackerel-plugin-repo1_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
			url,
			"Download URL is made correctly",
		)
	}
}

func TestInstallTargetGetOwnerAndRepo(t *testing.T) {
	{
		// it already has owner and repo
		it := &installTarget{
			owner: "owner1",
			repo:  "check-repo1",
		}
		owner, repo, err := it.getOwnerAndRepo()
		assert.Equal(t, "owner1", owner)
		assert.Equal(t, "check-repo1", repo)
		assert.NoError(t, err, "getOwnerAndRepo is finished successfully")
	}

	{
		// plugin def is not found in registry
		ts := httptest.NewServer(http.NotFoundHandler())
		defer ts.Close()

		it := &installTarget{
			pluginName:   "mackerel-plugin-not-found",
			rawGithubURL: ts.URL,
		}
		owner, repo, err := it.getOwnerAndRepo()
		assert.Equal(t, "", owner)
		assert.Equal(t, "", repo)
		assert.Error(t, err, "getOwnerAndRepo is failed because plugin def is not found")
		assert.Contains(t, err.Error(), "http response not OK. code: 404,", "Returns correct err")
	}

	{
		// plugin def is invalid json
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, `{"invalid" "jso"`)
		}))
		defer ts.Close()

		it := &installTarget{
			pluginName:   "mackerel-plugin-invalid-json",
			rawGithubURL: ts.URL,
		}
		owner, repo, err := it.getOwnerAndRepo()
		assert.Equal(t, "", owner)
		assert.Equal(t, "", repo)
		assert.Error(t, err, "getOwnerAndRepo is failed because plugin def is invalid json")
		assert.IsType(t, new(json.SyntaxError), err, "error type is syntax error")
	}

	{
		// a source field of plugin def is invalid
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, `{"description": "description", "source": "owner1"}`)
		}))
		defer ts.Close()

		it := &installTarget{
			pluginName:   "mackerel-plugin-invalid-source",
			rawGithubURL: ts.URL,
		}
		owner, repo, err := it.getOwnerAndRepo()
		assert.Equal(t, "", owner)
		assert.Equal(t, "", repo)
		assert.Error(t, err, "getOwnerAndRepo is failed because plugin def has invalid source")
		assert.Equal(t, err.Error(), "source definition is invalid", "Returns correct error")
	}

	{
		// get owner and repo correctly from registry
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/mackerelio/plugin-registry/master/plugins/mackerel-plugin-sample.json" {
				fmt.Fprint(w, `{"description": "Sample mackerel plugin", "source": "mackerelio/mackerel-plugin-sample"}`)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer ts.Close()

		it := &installTarget{
			pluginName:   "mackerel-plugin-sample",
			rawGithubURL: ts.URL,
		}
		owner, repo, err := it.getOwnerAndRepo()
		assert.Equal(t, "mackerelio", owner)
		assert.Equal(t, "mackerel-plugin-sample", repo)
		assert.NoError(t, err, "getOwnerAndRepo finished successfully")

		assert.Equal(t, "mackerelio", it.owner, "owner is cached")
		assert.Equal(t, "mackerel-plugin-sample", it.repo, "repo is cached")
	}
}

func TestInstallTargetGetReleaseTag(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/owner1/repo1/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"tag_name": "v0.5.1"}`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	{
		// it already has releaseTag
		it := &installTarget{releaseTag: "v0.1.2"}
		releaseTag, err := it.getReleaseTag("owner", "repo")
		assert.NoError(t, err, "getReleaseTag is successful")
		assert.Equal(t, "v0.1.2", releaseTag, "Returns correct releaseTag")
	}

	{
		// Specified owner and repo is not found
		it := &installTarget{githubAPIURL: ts.URL}
		releaseTag, err := it.getReleaseTag("owner1", "not-found-repo")
		assert.Error(t, err, "Returns err if the repository is not found")
		assert.Equal(t, "", releaseTag, "Returns empty string")
	}

	{
		// Get latest releaseTag correctly
		it := &installTarget{githubAPIURL: ts.URL}
		releaseTag, err := it.getReleaseTag("owner1", "repo1")
		assert.NoError(t, err, "getReleaseTag is successful")
		assert.Equal(t, "v0.5.1", releaseTag, "releaseTag is fetched correctly from api")
	}
}

func TestInstallTargetGetRawGithubURL(t *testing.T) {
	it := &installTarget{}
	assert.Equal(t, "https://raw.githubusercontent.com", it.getRawGithubURL(), "Returns default URL")

	it = &installTarget{rawGithubURL: "https://example.com"}
	assert.Equal(t, "https://example.com", it.getRawGithubURL(), "Returns customized URL")
}

func TestInstallTargetGetGithubAPIURL(t *testing.T) {
	it := &installTarget{}
	assert.Equal(t, "https://api.github.com/", it.getGithubAPIURL().String(), "Returns default URL")

	it = &installTarget{githubAPIURL: "https://api.example.com"}
	assert.Equal(t, "https://api.example.com/", it.getGithubAPIURL().String(), "Returns customized URL")
}
