package provider

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// Gzip provider
type Gzip struct {
	Path       string   // Path of the Gzip file
	file       *os.File // openned gzip file
	gzipReader *gzip.Reader
	tarReader  *tar.Reader
}

/*
func extractGzip(dest string, gzipPath string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractGzip: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return fmt.Errorf("ExtractGzip: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return fmt.Errorf("ExtractGzip: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractGzip: Copy() failed: %s", err.Error())
			}
			outFile.Close()
		}

	}
}
*/

// extractTarFromGzip extracts tar file inside a gzip file
func extractTarFromGzip(dest string, gzipPath string) error {
	inputFile, err := os.Open(gzipPath)
	if err != nil {
		return err
	}
	gzipReader, err := gzip.NewReader(inputFile)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	_, err = io.Copy(outputFile, gzipReader)
	if err != nil {
		return err
	}

	return nil
}

// Open opens the provider
func (c *Gzip) Open() (err error) {
	c.file, err = os.Open(c.Path)
	if err != nil {
		// TODO error management
		return ErrProviderUnavaiable
	}
	c.gzipReader, err = gzip.NewReader(c.file)
	fmt.Println(err)
	if err != nil {
		c.Close()
		return
	}
	c.tarReader = tar.NewReader(c.gzipReader)
	return nil
}

// Close closes the provider
func (c *Gzip) Close() error {
	if c.file == nil {
		return nil
	}
	err := c.file.Close()
	c.file = nil
	c.gzipReader = nil
	c.tarReader = nil
	return err
}

// GetLatestVersion gets the latest version
func (c *Gzip) GetLatestVersion() (string, error) {
	return "1.0", nil
}

// Walk walks all the files provided
func (c *Gzip) Walk(walkFn WalkFunc) error {
	if c.file == nil {
		return errors.New("nil file")
	}
	//c.gzipReader.Close()

	for true {
		header, err := c.tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Walk Gzip: Next() failed: %s", err.Error())
		}
		err = walkFn(&FileInfo{
			Path: header.Name,
			Mode: os.FileMode(header.Mode),
		})

	}
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *Gzip) Retrieve(src string, dest string) error {
	return nil
}
