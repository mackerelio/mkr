package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		url, err := it.makeDownloadURL(t.Context())
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
		url, err := it.makeDownloadURL(t.Context())
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
		_, err = it.makeDownloadURL(t.Context())
		assert.Error(t, err, "makeDownloadURL is failed")
	}

	{
		// Make download URL for `<owner>/<repo>` (latest release)
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/owner1/check-repo1/releases/latest", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"tag_name": "1.01"}`)
		})
		apiGithubServer := httptest.NewServer(mux)
		defer apiGithubServer.Close()

		it := &installTarget{
			owner:        "owner1",
			repo:         "check-repo1",
			apiGithubURL: apiGithubServer.URL,
			githubURL:    apiGithubServer.URL,
		}
		url, err := it.makeDownloadURL(t.Context())
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
			apiGithubURL: apiGithubServer.URL,
			githubURL:    apiGithubServer.URL,
		}
		_, err = it.makeDownloadURL(t.Context())
		assert.Error(t, err, "makeDownloadURL is failed")
	}

	{
		// Make download URL for `<pluginName>`
		muxAPI := http.NewServeMux()
		muxAPI.HandleFunc("/repos/owner1/mackerel-plugin-repo1/releases/latest", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"tag_name": "release/v0.5.1"}`)
		})
		apiGithubServer := httptest.NewServer(muxAPI)
		defer apiGithubServer.Close()

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
			apiGithubURL: apiGithubServer.URL,
			rawGithubURL: rawGithubServer.URL,
			githubURL:    rawGithubServer.URL,
		}

		url, err := it.makeDownloadURL(t.Context())
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
		owner, repo, err := it.getOwnerAndRepo(t.Context())
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
		owner, repo, err := it.getOwnerAndRepo(t.Context())
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
		owner, repo, err := it.getOwnerAndRepo(t.Context())
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
		owner, repo, err := it.getOwnerAndRepo(t.Context())
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
		owner, repo, err := it.getOwnerAndRepo(t.Context())
		assert.Equal(t, "mackerelio", owner)
		assert.Equal(t, "mackerel-plugin-sample", repo)
		assert.NoError(t, err, "getOwnerAndRepo finished successfully")

		assert.Equal(t, "mackerelio", it.owner, "owner is cached")
		assert.Equal(t, "mackerel-plugin-sample", it.repo, "repo is cached")
	}
}

func TestIsReleasesLatestURL(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output bool
	}{
		{
			Name:   "Valid releases/latest URL",
			Input:  "https://github.com/owner/repo/releases/latest",
			Output: true,
		},
		{
			Name:   "Invalid releases/latest URL",
			Input:  "https://github.com/owner/repo/releases/tag/v1.2.3",
			Output: false,
		},
		{
			Name:   "Invalid URL",
			Input:  "https://example.com/owner/repo/releases/latest",
			Output: false,
		},
		{
			Name:   "Invalid host with the github URL as its prefix",
			Input:  "https://github.com.example.com/owner/repo/releases/latest",
			Output: false,
		},
		{
			Name:   "Different scheme",
			Input:  "http://github.com/owner/repo/releases/latest",
			Output: false,
		},
		{
			Name:   "Extra segment before the owner",
			Input:  "https://github.com/extra/owner/repo/releases/latest",
			Output: false,
		},
		{
			Name:   "Path ending with a slash",
			Input:  "https://github.com/owner/repo/releases/latest/",
			Output: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			it := &installTarget{githubURL: "https://github.com"}
			u, _ := url.Parse(tc.Input)
			assert.Equal(t, tc.Output, it.isReleasesLatestURL(u))
		})
	}
}

func TestExtractTagFromReleasesURL(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output string
		Error  bool
	}{
		{
			Name:   "Valid releases/tag URL",
			Input:  "https://github.com/owner/repo/releases/tag/v1.2.3",
			Output: "v1.2.3",
		},
		{
			Name:   "Release tag containing an escaped slash",
			Input:  "https://github.com/owner/repo/releases/tag/release%2Fv0.5.1",
			Output: "release/v0.5.1",
		},
		{
			Name:   "Path ending with a slash",
			Input:  "https://github.com/owner/repo/releases/tag/v1.2.3/",
			Output: "v1.2.3",
		},
		{
			Name:  "Extra segment after the release tag",
			Input: "https://github.com/owner/repo/releases/tag/v1.2.3/extra",
			Error: true,
		},
		{
			Name:  "releases/tag in an unexpected position",
			Input: "https://github.com/releases/tag/owner/repo",
			Error: true,
		},
		{
			Name:  "releases/latest URL",
			Input: "https://github.com/owner/repo/releases/latest",
			Error: true,
		},
		{
			Name:  "Too few segments",
			Input: "https://github.com/owner/repo",
			Error: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			u, err := url.Parse(tc.Input)
			assert.NoError(t, err)
			tag, err := extractTagFromReleasesURL(u)
			if tc.Error {
				assert.Error(t, err, "returns an error for an unexpected URL")
				assert.Equal(t, "", tag, "returns an empty string")
				return
			}
			assert.NoError(t, err, "extractTagFromReleasesURL is successful")
			assert.Equal(t, tc.Output, tag, "release tag is extracted correctly")
		})
	}
}

func TestInstallTargetgetTagFromReleasesURL(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/owner1/repo1/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://github.com/owner1/repo1/releases/tag/v0.5.1")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	mux.HandleFunc("/owner2/repo2/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://github.com/other/path")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mux.HandleFunc("/old-owner/repo3/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", fmt.Sprintf("%s/new-owner/repo3/releases/latest", ts.URL))
		w.WriteHeader(http.StatusMovedPermanently)
	})
	mux.HandleFunc("/new-owner/repo3/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://github.com/new-owner/repo3/releases/tag/v1.0.0")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	mux.HandleFunc("/owner3/repo4/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", fmt.Sprintf("%s/owner3/repo4/releases/latest", ts.URL))
		w.WriteHeader(http.StatusMovedPermanently)
	})

	{
		// it already has releaseTag
		it := &installTarget{releaseTag: "v0.1.2"}
		releaseTag, err := it.getTagFromReleasesURL(t.Context(), "owner", "repo")
		assert.NoError(t, err, "getReleaseTag is successful")
		assert.Equal(t, "v0.1.2", releaseTag, "Returns correct releaseTag")
	}

	{
		// Specified owner and repo is not found
		it := &installTarget{githubURL: ts.URL}
		releaseTag, err := it.getTagFromReleasesURL(t.Context(), "owner1", "not-found-repo")
		assert.Error(t, err, "Returns err if the repository is not found")
		assert.Equal(t, "", releaseTag, "Returns empty string")
	}

	{
		// Get latest releaseTag correctly
		it := &installTarget{githubURL: ts.URL}
		releaseTag, err := it.getTagFromReleasesURL(t.Context(), "owner1", "repo1")
		assert.NoError(t, err, "getReleaseTag is successful")
		assert.Equal(t, "v0.5.1", releaseTag, "releaseTag is fetched correctly from api")
	}

	{
		// Specified owner and repo is not found (wrong redirect url)
		it := &installTarget{githubURL: ts.URL}
		releaseTag, err := it.getTagFromReleasesURL(t.Context(), "owner2", "repo2")
		assert.Error(t, err, "Returns err if the repository is not found")
		assert.Equal(t, "", releaseTag, "Returns empty string")
	}

	{
		// Repository transfer: follow redirect to new repository's releases/latest
		it := &installTarget{githubURL: ts.URL}
		releaseTag, err := it.getTagFromReleasesURL(t.Context(), "old-owner", "repo3")
		assert.NoError(t, err, "getReleaseTag is successful")
		assert.Equal(t, "v1.0.0", releaseTag, "releaseTag is fetched from transferred repository")
	}

	{
		// Too many redirects
		it := &installTarget{githubURL: ts.URL}
		releaseTag, err := it.getTagFromReleasesURL(t.Context(), "owner3", "repo4")
		assert.Error(t, err, "Returns err if too many redirects")
		assert.Equal(t, "", releaseTag, "Returns empty string")
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
		it := &installTarget{apiGithubURL: ts.URL}
		releaseTag, err := it.getReleaseTag("owner1", "not-found-repo")
		assert.Error(t, err, "Returns err if the repository is not found")
		assert.Equal(t, "", releaseTag, "Returns empty string")
	}

	{
		// Get latest releaseTag correctly
		it := &installTarget{apiGithubURL: ts.URL}
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

func TestInstallTargetGetAPIGithubURL(t *testing.T) {
	it := &installTarget{}
	assert.Equal(t, "https://api.github.com/", it.getAPIGithubURL().String(), "Returns default URL")

	it = &installTarget{apiGithubURL: "https://api.example.com"}
	assert.Equal(t, "https://api.example.com/", it.getAPIGithubURL().String(), "Returns customized URL")
}
