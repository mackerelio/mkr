package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
