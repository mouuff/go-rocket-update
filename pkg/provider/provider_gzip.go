package provider

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Gzip provider
type Gzip struct {
	Path          string // Path of the Gzip file
	tmpDir        string
	localProvider *Local
}

// extractGzip extracts gzip file to a folder
func extractGzip(tarball, dest string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dest, header.Name)
		info := header.FileInfo()
		if header.Typeflag == tar.TypeDir {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
		} else if header.Typeflag == tar.TypeReg {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tarReader)
			file.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Open opens the provider
func (c *Gzip) Open() (err error) {
	if c.tmpDir != "" {
		// If Open() has already been called we just ignore
		return nil
	}
	c.tmpDir, err = fileio.TempDir()
	if err != nil {
		return
	}
	err = extractGzip(c.Path, c.tmpDir)
	if err != nil {
		return err
	}
	c.localProvider = &Local{
		Path: c.tmpDir,
	}
	return nil
}

// Close closes the provider
func (c *Gzip) Close() (err error) {
	c.localProvider.Close()
	c.localProvider = nil
	if c.tmpDir != "" {
		err = os.RemoveAll(c.tmpDir)
		c.tmpDir = ""
		return
	}
	return
}

// GetLatestVersion gets the latest version
func (c *Gzip) GetLatestVersion() (string, error) {
	return GetLatestVersionFromPath(c.Path)
}

// Walk walks all the files provided
func (c *Gzip) Walk(walkFn WalkFunc) error {
	if c.localProvider == nil {
		return errors.New("nil c.localProvider")
	}
	return c.localProvider.Walk(walkFn)
}

// Retrieve file relative to "provider" to destination
func (c *Gzip) Retrieve(src string, dest string) error {
	if c.localProvider == nil {
		return errors.New("nil c.localProvider")
	}
	return c.localProvider.Retrieve(src, dest)
}
