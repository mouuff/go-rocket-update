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

	filesCount := 0
	err = p.Walk(func(info *provider.FileInfo) error {
		destPath := filepath.Join(tmpDir, info.Path)
		if info.IsDir {
			os.MkdirAll(destPath, os.ModePerm)
		} else {
			filesCount += 1
			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			err = p.Retrieve(info.Path, destPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if filesCount <= 0 {
		return fmt.Errorf("filesCount <= 0")
	}

	err = p.Walk(func(info *provider.FileInfo) error {
		destPath := filepath.Join(tmpDir, info.Path)
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
