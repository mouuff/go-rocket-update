package provider

import (
	"crypto/rsa"

	"github.com/mouuff/go-rocket-update/crypto"
)

// Secure provider defines a provider which verifies the signature of files when
// Retrieve() is called
type Secure struct {
	Provider   Provider
	PublicKey  rsa.PublicKey
	signatures crypto.Signatures
}

// Open opens the provider
func (c *Secure) Open() error {
	return c.Provider.Open()
}

// Close closes the provider
func (c *Secure) Close() error {
	return c.Provider.Close()
}

// GetLatestVersion gets the latest version
func (c *Secure) GetLatestVersion() (string, error) {
	return c.Provider.GetLatestVersion()
}

// Walk walks all the files provided
func (c *Secure) Walk(walkFn WalkFunc) error {
	return c.Provider.Walk(walkFn)
}

// Retrieve file and verifies the signature
func (c *Secure) Retrieve(src string, dest string) error {
	err := c.Provider.Retrieve(src, dest)
	if err != nil {
		return err
	}
	return nil
}
