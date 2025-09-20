package provider_test

import (
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
		return fmt.Errorf("bad version: %s", version)
	}
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tmpDest := filepath.Join(tmpDir, "tmpDest")
	err = p.Retrieve("thisfiledoesnotexists", tmpDest)
	if err == nil {
		return fmt.Errorf("provider.Retrieve() should return an error when source file does not exists")
	}
	if fileio.FileExists(tmpDest) {
		return fmt.Errorf("provider.Retrieve() should not create destination file when source file does not exists")
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
			return fmt.Errorf("file %s should exists", destPath)
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
		return fmt.Errorf("Walk() cancelled")
	})
	if err == nil {
		return fmt.Errorf("Walk() should return the error of walkFunc")
	}
	if count > 1 {
		return fmt.Errorf("Walk() should have stopped on error")
	}
	return nil
}

// ProviderTestUnavailable tests the expected behavior of a provider when it is not available
func ProviderTestUnavailable(p provider.Provider) error {
	if err := p.Open(); err == nil {
		return fmt.Errorf("Open() should return an error when provider is not available")
	}
	walkCount := 0
	err := p.Walk(func(info *provider.FileInfo) error {
		walkCount += 1
		return nil
	})

	if err == nil {
		return fmt.Errorf("Walk() should return an error when provider is not available")
	}

	if walkCount > 0 {
		return fmt.Errorf("Walk() should not call WalkFunc when provider is not available")
	}

	defer p.Close()
	_, err = p.GetLatestVersion()
	if err == nil {
		return fmt.Errorf("GetLatestVersion() should return an error when provider is not available")
	}
	if err = p.Close(); err != nil {
		return fmt.Errorf("Close() should not return an error if provider is not Open()")
	}
	return nil
}
