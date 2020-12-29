package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
)

// ChecksumFileSHA256 calculate the sha256 checksum of a file
func ChecksumFileSHA256(path string) ([]byte, error) {
	f, err := os.Open(path)
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

// GeneratePrivateKey generate a random private key.
func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

// GetFileSignature signs a file using the given private key
// returns the signature in a hex string
func GetFileSignature(priv *rsa.PrivateKey, path string) ([]byte, error) {
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

// VerifyFileSignature verifies the signature of a file
func VerifyFileSignature(pub *rsa.PublicKey, signature []byte, path string) error {
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

// ExportPrivateKeyAsPem exports the private key to Pem
func ExportPrivateKeyAsPem(privkey *rsa.PrivateKey) []byte {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return privkeyPem
}

// ParsePemPrivateKey parses the pem private key
func ParsePemPrivateKey(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// ExportPublicKeyAsPem exports the public key as Pem
func ExportPublicKeyAsPem(pubkey *rsa.PublicKey) ([]byte, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)
	return pubkeyPem, nil
}

// ParsePemPublicKey parse the pem public key
func ParsePemPublicKey(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}
	return nil, errors.New("Key type is not RSA")
}
