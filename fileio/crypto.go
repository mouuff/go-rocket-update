package fileio

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
)

// RandomPrivateKey generate a random private key.
func RandomPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

// GetSignature signs a file using the given private key
// returns the signature in a hex string
func GetSignature(priv *rsa.PrivateKey, path string) (string, error) {
	hash, err := ChecksumFile(path)
	if err != nil {
		return "", err
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies the signature of a file
func VerifySignature(pub *rsa.PublicKey, hexSignature string, path string) error {
	hash, err := ChecksumFile(path)
	if err != nil {
		return err
	}
	signature, err := hex.DecodeString(hexSignature)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash, signature)
	if err != nil {
		return err
	}
	return nil
}
