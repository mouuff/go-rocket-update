package crypto

import (
	"crypto/rsa"
	"errors"
	"os"
	"path/filepath"
)

// Signatures stores all the signatures of the files within a folder
type Signatures struct {
	Version       string
	SignaturesMap map[string][]byte
}

// GetFolderSignatures computes the signatures for a folder using the given private key
func GetFolderSignatures(priv *rsa.PrivateKey, root string) (*Signatures, error) {
	s := &Signatures{
		Version:       "1",
		SignaturesMap: map[string][]byte{},
	}
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.Mode().IsRegular() {
			signature, err := GetFileSignature(priv, filePath)
			if err != nil {
				return err
			}
			relpath, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			s.Add(relpath, signature)
		}
		return nil
	})
	return s, err
}

// VerifyFolder verifies all the files signatures of a folder
// returns list of unverified files
func (s *Signatures) VerifyFolder(pub *rsa.PublicKey, root string) ([]string, error) {
	unverifiedFiles := []string{}
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.Mode().IsRegular() {
			relpath, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			err = s.Verify(pub, relpath, filePath)
			if err != nil {
				unverifiedFiles = append(unverifiedFiles, filePath)
				return nil
			}
		}
		return nil
	})
	return unverifiedFiles, err
}

// Verify verifies the signature of a file
// full path must be the full path to the file
// relpath must be relative to the root
func (s *Signatures) Verify(pub *rsa.PublicKey, relpath string, fullpath string) error {
	// Why pass 'fullpath' instead of root? Because later on are going to retrieve files using a "Provider"
	// When we retrieve files using a "Provider", the file can be downloaded anywhere unrelated to the root
	signature, err := s.Get(relpath)
	if err != nil {
		return err
	}
	err = VerifyFileSignature(pub, signature, fullpath)
	if err != nil {
		return err
	}
	return err
}

// Add adds a signature of a file
// relpath must be a relative path from the root of the folder
func (s *Signatures) Add(relpath string, signature []byte) {
	s.SignaturesMap[filepath.ToSlash(relpath)] = signature
}

// Get gets a signature of a file given a relative path
func (s *Signatures) Get(relpath string) ([]byte, error) {
	if val, ok := s.SignaturesMap[filepath.ToSlash(relpath)]; ok {
		return val, nil
	}
	return nil, errors.New("Signature for file not found")
}
