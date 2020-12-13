package provider_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/easy-update/fileio"
	provider "github.com/mouuff/easy-update/provider"
)

func testProvider(p provider.Provider) error {
	tmpDir, err := ioutil.TempDir("", "testProvider")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = p.Walk(func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", filePath, err)
			return err
		}
		destPath := filepath.Join(tmpDir, filePath)
		if info.IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			err = p.Retrieve(filePath, destPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = p.Walk(func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", filePath, err)
			return err
		}
		destPath := filepath.Join(tmpDir, filePath)
		if !fileio.FileExists(destPath) {
			return fmt.Errorf("File %s should exists", destPath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func TestProviderLocal(t *testing.T) {
	p := provider.NewProviderLocal(filepath.Join("testdata", "Allum1"))
	if err := p.Open(); err != nil {
		t.Error(err)
	}
	defer p.Close()

	err := testProvider(p)
	if err != nil {
		t.Error(err)
	}

	tmpDir, err := ioutil.TempDir("", "TestProviderLocal")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	destPath := filepath.Join(tmpDir, "test.txt")
	err = p.Retrieve(filepath.Join("subfolder", "testfile.txt"), destPath)
	if err != nil {
		t.Error(err)
	}
	equals, err := fileio.CompareFileChecksum(destPath, filepath.Join("testdata", "Allum1", "subfolder", "testfile.txt"))
	if err != nil {
		t.Error(err)
	}
	if equals == false {
		t.Error("Files should be equals")
	}

}
