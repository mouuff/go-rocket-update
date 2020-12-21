package provider_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mouuff/easy-update/fileio"
	provider "github.com/mouuff/easy-update/provider"
)

func ProviderTestWalkAndRetrieve(p provider.Provider) error {
	tmpDir, err := ioutil.TempDir("", "testProvider")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = p.Walk(func(filePath string, isDir bool) error {
		destPath := filepath.Join(tmpDir, filePath)
		if isDir {
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

	err = p.Walk(func(filePath string, isDir bool) error {
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
