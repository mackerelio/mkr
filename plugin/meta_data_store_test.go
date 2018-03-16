package plugin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetaData(t *testing.T) {
	tmpd := tempd(t)
	defer os.RemoveAll(tmpd)

	it, err := newInstallTargetFromString("mackerelio/mackerel-plugin-sample@v1.0.1")
	assert.Nil(t, err, "error does not occur while newInstallTargetFromString")

	meta, err := newMetaDataStore(tmpd, it)
	assert.Nil(t, err, "error does not occur while newMetaDataStore")

	v1, err := meta.load("foo")
	assert.Nil(t, err, "error does not occur while meta.load")
	assert.Equal(t, v1, "", "load uninitialized value must be empty string")

	err = meta.store("foo", "bar")
	assert.Nil(t, err, "error does not occur while meta.store")

	v2, err := meta.load("foo")
	assert.Nil(t, err, "error does not occur while meta.load")
	assert.Equal(t, v2, "bar", "load successful")

	err = meta.store("foo", "baz")
	assert.Nil(t, err, "error does not occur while meta.store")

	v3, err := meta.load("foo")
	assert.Nil(t, err, "error does not occur while meta.load")
	assert.Equal(t, v3, "baz", "load successful")

	_, err = meta.load("")
	assert.Error(t, err, "error occured while meta.load with empty key")
}
