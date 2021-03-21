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

// ProviderTestWalkAndRetrieve tests the expected behavior of a provider
func ProviderTestWalkAndRetrieve(p provider.AccessProvider) error {
	version, err := p.GetLatestVersion()
	if err != nil {
		return err
	}
	if len(version) < 1 { // TODO idea check version format?
		return errors.New("Bad version: " + version)
	}
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tmpDest := filepath.Join(tmpDir, "tmpDest")
	err = p.Retrieve("thisfiledoesnotexists", tmpDest)
	if err == nil {
		return errors.New("provider.Retrieve() should return an error when source file does not exists")
	}
	if fileio.FileExists(tmpDest) {
		return errors.New("provider.Retrieve() should not create destination file when source file does not exists")
	}

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

// ProviderTestUnavaiable tests the expected behavior of a provider when it is not avaiable
func ProviderTestUnavaiable(p provider.Provider) error {
	if err := p.Open(); err == nil {
		return errors.New("Open() should return an error when provider is not avaiable")
	}
	_, err := p.GetLatestVersion()
	if err == nil {
		return errors.New("GetLatestVersion() should return an error when provider is not avaiable")
	}
	if err = p.Close(); err != nil {
		return errors.New("Close() should not return an error if provider is not Open()")
	}
	return nil
}
