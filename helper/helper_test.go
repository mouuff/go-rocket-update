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

	srcChecksum, err := helper.ChecksumFile(src)
	if err != nil {
		return err
	}
	destChecksum, err := helper.ChecksumFile(dest)
	if err != nil {
		return err
	}
	if destChecksum != srcChecksum {
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
