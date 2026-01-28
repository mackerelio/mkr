package plugin

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type metaDataStore struct {
	dir           string
	installTarget *installTarget
}

var errDisaleMetaDataStore = errors.New("MetaData disabled. could not detect owner/repo")

func newMetaDataStore(ctx context.Context, pluginDir string, target *installTarget) (*metaDataStore, error) {
	owner, repo, err := target.getOwnerAndRepo(ctx)
	if err != nil {
		return nil, errDisaleMetaDataStore
	}
	dir := filepath.Join(pluginDir, "meta", owner, repo)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &metaDataStore{
		dir:           dir,
		installTarget: target,
	}, nil
}

func (m *metaDataStore) load(key string) (string, error) {
	b, err := os.ReadFile(filepath.Join(m.dir, key))
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m *metaDataStore) store(key, value string) error {
	return os.WriteFile(
		filepath.Join(m.dir, key),
		[]byte(value),
		0644,
	)
}
