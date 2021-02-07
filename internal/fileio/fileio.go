package fileio

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/kardianos/osext"
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

// FileExists checks if the file exists
func FileExists(src string) bool {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return false
	}
	return true
}

// ChecksumFile calculate the checksum of a file
// This is used only internally to compare files
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

// CompareFiles compares two files
// returns True if files are the same
func CompareFiles(fileA, fileB string) (bool, error) {
	fileAChecksum, err := ChecksumFile(fileA)
	if err != nil {
		return false, err
	}
	fileBChecksum, err := ChecksumFile(fileB)
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

// GetExecutable get the path to the current executable
func GetExecutable() (string, error) {
	execPath, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return execPath, nil
}
