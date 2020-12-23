package provider

import (
	"archive/zip"
	"io"
	"os"
)

// A Zip provider provides files that are inside a zip folder
type Zip struct {
	Path   string          // Path of the zip file
	reader *zip.ReadCloser // reader for the current zip file
}

// Open opens the provider
func (c *Zip) Open() error {
	_, err := os.Stat(c.Path)
	if os.IsNotExist(err) {
		return ErrProviderUnavaiable
	}
	c.reader, err = zip.OpenReader(c.Path)
	if err != nil {
		c.reader = nil
		return err
	}
	return nil
}

// Close closes the provider
func (c *Zip) Close() error {
	return c.reader.Close()
}

// GetLatestVersion gets the lastest version
func (c *Zip) GetLatestVersion() (string, error) {
	return "1.0", nil
}

// Walk walks all the files provided
func (c *Zip) Walk(walkFn WalkFunc) error {
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
func (c *Zip) findFileByPath(path string) *zip.File {
	for _, f := range c.reader.File {
		if f.Name == path {
			return f
		}
	}
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *Zip) Retrieve(src string, dest string) error {

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
