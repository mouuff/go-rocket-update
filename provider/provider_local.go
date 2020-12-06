package provider

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mouuff/easy-update/fileio"
)

type ProviderLocal struct {
	path string
}

func NewProviderLocal(path string) Provider {
	return &ProviderLocal{path: path}
}

func (c *ProviderLocal) Open() error {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return ErrProviderUnavaiable
	}
	return nil
}

func (c *ProviderLocal) Close() error {
	return nil
}

func (c *ProviderLocal) GetVersion() (string, error) {
	return "1.0", nil
}

func (c *ProviderLocal) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(c.path, func(filePath string, info os.FileInfo, walkErr error) error {
		relpath, err := filepath.Rel(c.path, filePath)
		if err != nil {
			return err
		}
		return walkFn(relpath, info, walkErr)
	})
}

func (c *ProviderLocal) Retrieve(src string, dest string) error {
	fullPath := path.Join(c.path, src)
	return fileio.CopyFile(fullPath, dest)
}
