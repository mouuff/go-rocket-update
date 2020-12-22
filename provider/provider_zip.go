package provider

import (
	"archive/zip"
	"os"
	"path/filepath"

	"github.com/mouuff/easy-update/fileio"
)

type providerZip struct {
	path   string
	reader *zip.ReadCloser
}

// NewProviderZip creates a new provider for local files
func NewProviderZip(path string) Provider {
	return &providerZip{path: path}
}

// Open opens the provider
func (c *providerZip) Open() error {
	_, err := os.Stat(c.path)
	if os.IsNotExist(err) {
		return ErrProviderUnavaiable
	}
	c.reader, err = zip.OpenReader(c.path)
	if err != nil {
		c.reader = nil
		return err
	}
	return nil
}

// Close closes the provider
func (c *providerZip) Close() error {
	return c.reader.Close()
}

// GetLatestVersion gets the lastest version
func (c *providerZip) GetLatestVersion() (string, error) {
	return "1.0", nil
}

// Walk walks all the files provided
func (c *providerZip) Walk(walkFn WalkFunc) error {
	for _, f := range c.reader.File {
		walkFn(&FileInfo{
			Path: f.Name,
			Mode: f.Mode(),
		})
	}
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *providerZip) Retrieve(src string, dest string) error {
	fullPath := filepath.Join(c.path, src)
	return fileio.CopyFile(fullPath, dest)
}
