package provider

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mouuff/easy-update/fileio"
)

type providerLocal struct {
	path string
}

// NewProviderLocal creates a new provider for local files
func NewProviderLocal(path string) Provider {
	return &providerLocal{path: path}
}

// Open opens the provider
func (c *providerLocal) Open() error {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return ErrProviderUnavaiable
	}
	return nil
}

// Close closes the provider
func (c *providerLocal) Close() error {
	return nil
}

// GetLatestVersion gets the lastest version
func (c *providerLocal) GetLatestVersion() (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(c.path, "VERSION"))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Walk walks all the files provided
func (c *providerLocal) Walk(walkFn WalkFunc) error {
	return filepath.Walk(c.path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr == nil {
			relpath, err := filepath.Rel(c.path, filePath)
			if err != nil {
				return err
			}
			return walkFn(&FileInfo{
				Path:  relpath,
				IsDir: info.IsDir(),
				Mode:  info.Mode(),
			})
		}
		return nil
	})
}

// Retrieve file relative to "provider" to destination
func (c *providerLocal) Retrieve(src string, dest string) error {
	fullPath := filepath.Join(c.path, src)
	return fileio.CopyFile(fullPath, dest)
}
