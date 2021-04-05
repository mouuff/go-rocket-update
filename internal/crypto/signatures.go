package crypto

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
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
			relPath, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			s.Add(relPath, signature)
		}
		return nil
	})
	return s, err
}

// LoadSignaturesFromJSON loads signatures from a JSON file
func LoadSignaturesFromJSON(path string) (signatures *Signatures, err error) {
	signaturesJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	signatures = &Signatures{}
	err = json.Unmarshal(signaturesJSON, signatures)
	if err != nil {
		return
	}
	return
}

// WriteSignaturesToJSON writes signatures to dest as a JSON file
func WriteSignaturesToJSON(dest string, signatures *Signatures) error {
	signaturesJSON, err := json.Marshal(signatures)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dest, signaturesJSON, 0644)
	if err != nil {
		return err
	}
	return nil
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
			relPath, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			err = s.Verify(pub, relPath, filePath)
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
// relPath must be relative to the root
func (s *Signatures) Verify(pub *rsa.PublicKey, relPath string, fullpath string) error {
	// Why pass 'fullpath' instead of root? Because later on are going to retrieve files using a "Provider"
	// When we retrieve files using a "Provider", the file can be downloaded anywhere unrelated to the root
	signature, err := s.Get(relPath)
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
// relPath must be a relative path from the root of the folder
func (s *Signatures) Add(relPath string, signature []byte) {
	s.SignaturesMap[filepath.ToSlash(relPath)] = signature
}

// Get gets a signature of a file given a relative path
func (s *Signatures) Get(relPath string) ([]byte, error) {
	if val, ok := s.SignaturesMap[filepath.ToSlash(relPath)]; ok {
		return val, nil
	}
	return nil, errors.New("Signature for file not found")
}

// Remove removes a signature of a file given a relative path
func (s *Signatures) Remove(relPath string) {
	_, ok := s.SignaturesMap[filepath.ToSlash(relPath)]
	if ok {
		delete(s.SignaturesMap, filepath.ToSlash(relPath))
	}
}
