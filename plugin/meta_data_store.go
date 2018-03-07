package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type metaDataStore struct {
	dir           string
	installTarget *installTarget
}

func newMetaDataStore(pluginDir string, it *installTarget) (*metaDataStore, error) {
	dir := filepath.Join(pluginDir, "meta", it.owner, it.repo)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &metaDataStore{
		dir:           dir,
		installTarget: it,
	}, nil
}

func (m *metaDataStore) load(key string) (string, error) {
	f, err := os.OpenFile(
		filepath.Join(m.dir, key),
		os.O_RDONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func (m *metaDataStore) store(key, value string) error {
	return ioutil.WriteFile(
		filepath.Join(m.dir, key),
		[]byte(value),
		0644,
	)
}
