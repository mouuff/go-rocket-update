package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

func ChecksumFile(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
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
	fileAChecksum, err := ChecksumFile(fileA)
	if err != nil {
		return false, err
	}
	fileBChecksum, err := ChecksumFile(fileB)
	if err != nil {
		return false, err
	}
	if fileBChecksum != fileAChecksum {
		return false, fmt.Errorf("fileBChecksum: %s != fileAChecksum: %s", fileB, fileA)
	}
	return true, nil
}
