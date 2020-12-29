package crypto

import (
	"errors"
	"path/filepath"
)

// FolderSignature stores all the signatures of the files within a folder
type FolderSignature struct {
	Version    string
	Signatures map[string][]byte
}

// NewFolderSignature instanciates a new folder signature
func NewFolderSignature() *FolderSignature {
	return &FolderSignature{
		Version:    "1",
		Signatures: map[string][]byte{},
	}
}

// AddSignature adds a signature of a file
// relpath must be a relative path from the root of the folder
func (fs *FolderSignature) AddSignature(relpath string, signature []byte) {
	fs.Signatures[filepath.ToSlash(relpath)] = signature
}

// GetSignature gets a signature of a file given a relative path
func (fs *FolderSignature) GetSignature(relpath string) ([]byte, error) {
	if val, ok := fs.Signatures[filepath.ToSlash(relpath)]; ok {
		return val, nil
	}
	return nil, errors.New("Signature for file not found")
}
