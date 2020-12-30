package provider

import (
	"crypto/rsa"
	"os"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/internal/constant"
	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Secure provider defines a provider which verifies the signature of files when
// Retrieve() is called
type Secure struct {
	Provider   Provider
	PublicKey  *rsa.PublicKey
	signatures *crypto.Signatures
}

// Open the provider
func (c *Secure) Open() error {
	err := c.Provider.Open()
	if err != nil {
		return err
	}
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, constant.SignatureRelPath)
	err = c.Provider.Retrieve(constant.SignatureRelPath, tmpFile)
	if err != nil {
		// TODO defines error
		return err
	}
	c.signatures, err = crypto.LoadSignaturesFromJSON(tmpFile)
	if err != nil {
		// TODO defines error
		return err
	}
	return nil
}

// Close the provider
func (c *Secure) Close() error {
	return c.Provider.Close()
}

// GetLatestVersion gets the latest version
func (c *Secure) GetLatestVersion() (string, error) {
	return c.Provider.GetLatestVersion()
}

// Walk all the files provided
func (c *Secure) Walk(walkFn WalkFunc) error {
	return c.Provider.Walk(walkFn)
}

// Retrieve file and verifies the signature
func (c *Secure) Retrieve(src string, dest string) error {
	err := c.Provider.Retrieve(src, dest)
	if err != nil {
		return err
	}
	err = c.signatures.Verify(c.PublicKey, src, dest)
	if err != nil {
		os.Remove(dest)
		return err
	}
	return nil
}
