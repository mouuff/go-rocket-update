package provider

import (
	"path/filepath"
)

type ProviderLocal struct {
	path string
}

func NewProviderLocal(path string) Provider {
	return &ProviderLocal{path: path}
}

func (c *ProviderLocal) Open() error {
	return nil
}

func (c *ProviderLocal) Close() error {
	return nil
}

func (c *ProviderLocal) GetVersion() (string, error) {
	return "", nil
}

func (c *ProviderLocal) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(c.path, walkFn)
}

func (c *ProviderLocal) Retrieve(srcPath string, destPath string) error {
	return filepath.Rel(c.path, srcPath)
}
