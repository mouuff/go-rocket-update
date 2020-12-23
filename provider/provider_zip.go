package provider

import (
	"archive/zip"
	"io"
	"os"
)

// Zip provider
type Zip struct {
	Reader *zip.ReadCloser // reader for the current zip file
}

// NewZipProvider creates a new zip provider from a zip file
func NewZipProvider(Path string) (c *Zip, err error) {
	c = &Zip{}
	c.Reader, err = zip.OpenReader(Path)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Open opens the provider
func (c *Zip) Open() error {
	return nil
}

// Close closes the provider
func (c *Zip) Close() error {
	if c.Reader == nil {
		return nil
	}
	return c.Reader.Close()
}

// GetLatestVersion gets the lastest version
func (c *Zip) GetLatestVersion() (string, error) {
	return "1.0", nil
}

// Walk walks all the files provided
func (c *Zip) Walk(walkFn WalkFunc) error {
	for _, f := range c.Reader.File {
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
	for _, f := range c.Reader.File {
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
