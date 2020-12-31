package provider_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mouuff/go-rocket-update/internal/constant"
	"github.com/mouuff/go-rocket-update/internal/fileio"
	"github.com/mouuff/go-rocket-update/pkg/provider"
)

func ProviderTestWalkAndRetrieve(p provider.AccessProvider) error {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	filesCount := 0
	err = p.Walk(func(info *provider.FileInfo) error {
		destPath := filepath.Join(tmpDir, info.Path)
		if info.Mode.IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
		} else {
			if strings.Contains(info.Path, constant.SignatureRelPath) {
				return nil
			}
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
		if !fileio.FileExists(destPath) && !strings.Contains(info.Path, constant.SignatureRelPath) {
			return fmt.Errorf("File %s should exists", destPath)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Test to make sure the walk stops when walkFunc returns an error
	count := 0
	err = p.Walk(func(info *provider.FileInfo) error {
		count += 1
		return errors.New("Walk cancelled")
	})
	if err == nil {
		return errors.New("Walk should return the error of walkFunc")
	}
	if count > 1 {
		return errors.New("Walk should have stopped on error")
	}
	return nil
}
