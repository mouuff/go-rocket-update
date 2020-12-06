package helper_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/mouuff/easy-update/helper"
)

func copyAndChecksumFile(src string) error {
	dir, err := ioutil.TempDir("", "copyAndChecksumFile")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	dest := path.Join(dir, "dest.txt")
	err = helper.CopyFile(src, dest)
	if err != nil {
		return err
	}

	equals, err := helper.CompareFileChecksum(src, dest)
	if err != nil {
		return err
	}
	if equals == false {
		return fmt.Errorf("destChecksum: %s != srcChecksum: %s", dest, src)
	}
	return nil
}

func TestCopyFile(t *testing.T) {
	err := copyAndChecksumFile(path.Join("testdata", "small.exe"))
	if err != nil {
		t.Error(err)
	}
	err = copyAndChecksumFile(path.Join("testdata", "TempleOS.ISO"))
	if err != nil {
		t.Error(err)
	}
	err = copyAndChecksumFile(path.Join("testdata", "empty.txt"))
	if err != nil {
		t.Error(err)
	}
}

func verifyChecksumFile(src, expectedChecksum string) error {
	checksum, err := helper.ChecksumFile(src)
	if err != nil {
		return err
	}
	if checksum != expectedChecksum {
		return fmt.Errorf("checksum: %s != expectedChecksum: %s", checksum, expectedChecksum)
	}
	return nil
}

func TestChecksum(t *testing.T) {
	err := verifyChecksumFile(path.Join("testdata", "empty.txt"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Error(err)
	}
	err = verifyChecksumFile(path.Join("testdata", "TempleOS.ISO"),
		"5d0fc944e5d89c155c0fc17c148646715bc1db6fa5750c0b913772cfec19ba26")
	if err != nil {
		t.Error(err)
	}

	// Test wrong checksum:
	err = verifyChecksumFile(path.Join("testdata", "TempleOS.ISO"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err == nil {
		t.Error(fmt.Errorf("verifyChecksumFile returned not nil on a bad checksum"))
	}

}
