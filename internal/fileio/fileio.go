package fileio

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	return os.MkdirTemp("", "rocket-updater")
}

// GetExecutable get the path to the current executable
func GetExecutable() (string, error) {
	return os.Executable()
}

// SafeJoin joins name onto dest and guarantees the result stays within dest.
// It rejects archive entries that use ".." to escape the extraction directory
// (Zip Slip, CWE-22), so entry names taken from an untrusted archive cannot be
// used to write outside dest.
func SafeJoin(dest, name string) (string, error) {
	cleanDest := filepath.Clean(dest)
	path := filepath.Join(cleanDest, name)
	if path != cleanDest && !strings.HasPrefix(path, cleanDest+string(os.PathSeparator)) {
		return "", fmt.Errorf("fileio: illegal archive path (path traversal): %q", name)
	}
	return path, nil
}
