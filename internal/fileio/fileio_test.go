package fileio_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/fileio"
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

	equals, err := fileio.CompareFiles(src, dest)
	if err != nil {
		return err
	}
	if equals == false {
		return fmt.Errorf("destChecksum: %s != srcChecksum: %s", dest, src)
	}
	return nil
}

func TestCopyFile(t *testing.T) {
	err := copyAndChecksumFile(filepath.Join("testdata", "smallexe"))
	if err != nil {
		t.Fatal(err)
	}
	err = copyAndChecksumFile(filepath.Join("testdata", "TempleOS.ISO"))
	if err != nil {
		t.Fatal(err)
	}
	err = copyAndChecksumFile(filepath.Join("testdata", "empty.txt"))
	if err != nil {
		t.Fatal(err)
	}

	err = fileio.CopyFile(filepath.Join("testdata", "smallexe"), filepath.Join("pathdoesnotexists", "pathdoesnotexists"))
	if err == nil {
		t.Fatal("CopyFile should not work if path does not exists")
	}
	err = fileio.CopyFile(filepath.Join("pathdoesnotexists", "pathdoesnotexists"), filepath.Join("testdata", "smallexe"))
	if err == nil {
		t.Fatal("CopyFile should not work if path does not exists")
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
		t.Fatal(err)
	}
	err = verifyChecksumFile(filepath.Join("testdata", "TempleOS.ISO"),
		"5d0fc944e5d89c155c0fc17c148646715bc1db6fa5750c0b913772cfec19ba26")
	if err != nil {
		t.Fatal(err)
	}

	// Test wrong checksum:
	err = verifyChecksumFile(filepath.Join("testdata", "TempleOS.ISO"),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err == nil {
		t.Fatal(fmt.Errorf("verifyChecksumFile returned not nil on a bad checksum"))
	}
}

func TestCompareFileChecksum(t *testing.T) {
	fileA := filepath.Join("testdata", "TempleOS.ISO")
	fileB := filepath.Join("testdata", "smallexe")
	fileC := filepath.Join("testdata", "doesNotExists")

	equals, err := fileio.CompareFiles(fileA, fileA)
	if err != nil {
		t.Fatal(err)
	}
	if equals == false {
		t.Fatal("Should be equal")
	}

	equals, err = fileio.CompareFiles(fileA, fileB)
	if err != nil {
		t.Fatal(err)
	}
	if equals == true {
		t.Fatal("Should be unequal")
	}

	_, err = fileio.CompareFiles(fileA, fileC)
	if err == nil {
		t.Fatal("fileio.CompareFiles(fileA, fileC) should return an error")
	}
	_, err = fileio.CompareFiles(fileC, fileB)
	if err == nil {
		t.Fatal("fileio.CompareFiles(fileC, fileB) should return an error")
	}
}

func TestFileExists(t *testing.T) {
	fileA := filepath.Join("testdata", "TempleOS.ISO")
	fileB := filepath.Join("testdata", "doesNotExists")

	if !fileio.FileExists(fileA) {
		t.Error("fileio.FileExists(fileA) exists")
	}
	if fileio.FileExists(fileB) {
		t.Error("fileio.FileExists(fileA) does not exists")
	}
}
