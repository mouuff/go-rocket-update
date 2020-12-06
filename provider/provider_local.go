package provider

import (
	"path/filepath"

	"github.com/mouuff/easy-update/helper"
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
	return "1.0", nil
}

func (c *ProviderLocal) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(c.path, walkFn)
}

func (c *ProviderLocal) Retrieve(src string, dest string) error {
	fullPath, err := filepath.Rel(c.path, src)
	if err != nil {
		return err
	}
	return helper.CopyFile(fullPath, dest)
}
