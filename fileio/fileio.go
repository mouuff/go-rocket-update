package fileio

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
)

// CopyFile copies file contents from file source to file destination
func CopyFile(src string, dest string) error {

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = out.Sync()
	return err
}

// ChecksumFile calculate the sha256 checksum of a file
func ChecksumFile(src string) ([]byte, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hash.Sum(nil), nil
}

// ChecksumFileHex calculate the sha256 checksum of a file returns hex value
func ChecksumFileHex(src string) (string, error) {
	b, err := ChecksumFile(src)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// FileExists checks if the file exists
func FileExists(src string) bool {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return false
	}
	return true
}

// CompareFileChecksum compares two files checksums
func CompareFileChecksum(fileA, fileB string) (bool, error) {
	fileAChecksum, err := ChecksumFileHex(fileA)
	if err != nil {
		return false, err
	}
	fileBChecksum, err := ChecksumFileHex(fileB)
	if err != nil {
		return false, err
	}
	if fileBChecksum != fileAChecksum {
		return false, nil
	}
	return true, nil
}

// TempDir creates a new temporary directory
func TempDir() (string, error) {
	return ioutil.TempDir("", "rocket-updater")
}

// SignFile signs a file using the given private key
// returns the signature as hex
func SignFile(priv *rsa.PrivateKey, path string) (string, error) {
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

// VerifyFile verifies the signature of a file
func VerifyFile(priv *rsa.PublicKey, path string) (string, error) {

}
