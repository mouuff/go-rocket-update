package crypto

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
	fs.Signatures[relpath] = signature
}

func (fs *FolderSignature) GetSignature(relpath string) ([]byte, error) {
	return fs.Signatures[relpath], nil
}
