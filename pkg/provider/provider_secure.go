package provider

import (
	"crypto/rsa"
	"os"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/internal/constant"
	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Secure is a provider which uses another provider and verifies the
// signatures when Retrieve() is called
// If you pass a nil PublicKey, then it will try to load the PublicKeyPEM.
// PublicKeyPEM must be a public key in the PEM format
type Secure struct {
	BackendProvider Provider
	PublicKeyPEM    []byte
	PublicKey       *rsa.PublicKey
	signatures      *crypto.Signatures
}

// Open the provider
func (c *Secure) Open() (err error) {
	if c.PublicKey == nil {
		c.PublicKey, err = crypto.ParsePemPublicKey(c.PublicKeyPEM)
		if err != nil {
			return
		}
	}
	err = c.BackendProvider.Open()
	if err != nil {
		return
	}
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, constant.SignatureRelPath)
	err = c.BackendProvider.Retrieve(constant.SignatureRelPath, tmpFile)
	if err != nil {
		// TODO defines error
		return
	}
	c.signatures, err = crypto.LoadSignaturesFromJSON(tmpFile)
	if err != nil {
		// TODO defines error
		return
	}
	return nil
}

// Close the provider
func (c *Secure) Close() error {
	return c.BackendProvider.Close()
}

// GetLatestVersion gets the latest version
func (c *Secure) GetLatestVersion() (string, error) {
	return c.BackendProvider.GetLatestVersion()
}

// Walk all the files provided
func (c *Secure) Walk(walkFn WalkFunc) error {
	return c.BackendProvider.Walk(walkFn)
}

// Retrieve file and verifies the signature
func (c *Secure) Retrieve(src string, dest string) error {
	err := c.BackendProvider.Retrieve(src, dest)
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
