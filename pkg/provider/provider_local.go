package provider

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// A Local provider use a local directory to provide files
// This provider is mainly here for mocking and testing but it could be used on a shared network folder
type Local struct {
	Path string // Path of the folder
}

// Open opens the provider
func (c *Local) Open() error {
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		return ErrProviderUnavaiable
	}
	return nil
}

// Close closes the provider
func (c *Local) Close() error {
	return nil
}

// GetLatestVersion gets the latest version
func (c *Local) GetLatestVersion() (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(c.Path, "VERSION"))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Walk walks all the files provided
func (c *Local) Walk(walkFn WalkFunc) error {
	return filepath.Walk(c.Path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			// TODO log walkErr ?
			return nil // Ignore error and continue walking
		}
		relPath, err := filepath.Rel(c.Path, filePath)
		if err != nil {
			return err
		}
		return walkFn(&FileInfo{
			Path: relPath,
			Mode: info.Mode(),
		})
	})
}

// Retrieve file relative to "provider" to destination
func (c *Local) Retrieve(src string, dest string) error {
	fullPath := filepath.Join(c.Path, src)
	return fileio.CopyFile(fullPath, dest)
}
