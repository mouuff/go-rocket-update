package fileio_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/easy-update/fileio"
)

func copyAndChecksumFile(src string) error {
	dir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	dest := filepath.Join(dir, "dest.txt")
	err = fileio.CopyFile(src, dest)
	if err != nil {
		return err
	}

	equals, err := fileio.CompareFileChecksum(src, dest)
	if err != nil {
		return err
	}
	if equals == false {
		return fmt.Errorf("destChecksum: %s != srcChecksum: %s", dest, src)
	}
	return nil
}

func TestCopyFile(t *testing.T) {
	err := copyAndChecksumFile(filepath.Join("testdata", "small.exe"))
	if err != nil {
		t.Error(err)
	}
	err = copyAndChecksumFile(filepath.Join("testdata", "TempleOS.ISO"))
	if err != nil {
		t.Error(err)
	}
	err = copyAndChecksumFile(filepath.Join("testdata", "empty.txt"))
	if err != nil {
		t.Error(err)
	}
}

func verifyChecksumFile(src, expectedChecksum string) error {
	checksum, err := fileio.ChecksumFile(src)
	if err != nil {
		return err
	}
	if checksum != expectedChecksum {
		return fmt.Errorf("checksum: %s != expectedChecksum: %s", checksum, expectedChecksum)
	}
	return nil
}

func TestChecksum(t *testing.T) {
	err := verifyChecksumFile(filepath.Join("testdata", "empty.txt"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Error(err)
	}
	err = verifyChecksumFile(filepath.Join("testdata", "TempleOS.ISO"),
		"5d0fc944e5d89c155c0fc17c148646715bc1db6fa5750c0b913772cfec19ba26")
	if err != nil {
		t.Error(err)
	}

	// Test wrong checksum:
	err = verifyChecksumFile(filepath.Join("testdata", "TempleOS.ISO"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err == nil {
		t.Error(fmt.Errorf("verifyChecksumFile returned not nil on a bad checksum"))
	}
}

func TestCompareFileChecksum(t *testing.T) {
	fileA := filepath.Join("testdata", "TempleOS.ISO")
	fileB := filepath.Join("testdata", "small.exe")
	equals, err := fileio.CompareFileChecksum(fileA, fileA)
	if err != nil {
		t.Error(err)
	}
	if equals == false {
		t.Error("Should be equal")
	}

	equals, err = fileio.CompareFileChecksum(fileA, fileB)
	if err != nil {
		t.Error(err)
	}
	if equals == true {
		t.Error("Should be unequal")
	}
}
