package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"os"
)

// ChecksumFileSHA256 calculate the sha256 checksum of a file
func ChecksumFileSHA256(src string) ([]byte, error) {
	f, err := os.Open(src)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return []byte{}, err
	}

	return hash.Sum(nil), nil
}

// RandomPrivateKey generate a random private key.
func RandomPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

// GetSignature signs a file using the given private key
// returns the signature in a hex string
func GetSignature(priv *rsa.PrivateKey, path string) ([]byte, error) {
	hash, err := ChecksumFileSHA256(path)
	if err != nil {
		return nil, err
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// VerifySignature verifies the signature of a file
func VerifySignature(pub *rsa.PublicKey, signature []byte, path string) error {
	hash, err := ChecksumFileSHA256(path)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash, signature)
	if err != nil {
		return err
	}
	return nil
}
