package provider

import (
	"archive/zip"
	"io"
	"os"
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

// findFileByPath finds a file in the currend zip by the path
// returns nil if file does not exists
func (c *providerZip) findFileByPath(path string) *zip.File {
	for _, f := range c.reader.File {
		if f.Name == path {
			return f
		}
	}
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *providerZip) Retrieve(src string, dest string) error {

	zipFile := c.findFileByPath(src)
	if zipFile == nil {
		return ErrFileNotFound
	}

	inputFile, err := zipFile.Open()
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(
		dest,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		zipFile.Mode(),
	)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}
	return nil
}
