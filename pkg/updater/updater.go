package updater

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mouuff/go-rocket-update/internal/fileio"
	"github.com/mouuff/go-rocket-update/pkg/provider"
)

// Updater struct
type Updater struct {
	Provider   provider.Provider
	BinaryName string
	Version    string
}

// GetBinaryPatcher gets the binary patcher
// binaryCandidate can be empty if you only plan to rollback
func GetBinaryPatcher(binaryCandidate string) (*Patcher, error) {
	executable, err := fileio.GetExecutable()
	if err != nil {
		return nil, err
	}
	return &Patcher{
		DestinationPath: executable,
		SourcePath:      binaryCandidate,
		BackupPath:      executable + ".old",
		Mode:            0755,
		Verify:          nil, // TODO
	}, nil
}

// getBinaryName gets the name used to find the right binary
func (u *Updater) getBinaryName() string {
	return u.BinaryName + "_" + runtime.GOOS + "_" + runtime.GOARCH
}

// findBinaryPath finds the right binary using the provider
func (u *Updater) findBinaryPath() (string, error) {
	binaryPath := ""
	err := u.Provider.Walk(func(info *provider.FileInfo) error {
		if info.Mode.IsRegular() && strings.Contains(filepath.Base(info.Path), u.getBinaryName()) {
			binaryPath = info.Path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return binaryPath, nil
}

// updateExecutable updates the current executable with the new one
func (u *Updater) updateExecutable() (err error) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)
	binaryPath, err := u.findBinaryPath()
	if err != nil {
		return
	}
	binaryTmpPath := filepath.Join(tmpDir, filepath.Base(binaryPath))
	err = u.Provider.Retrieve(binaryPath, binaryTmpPath)
	if err != nil {
		return
	}

	patcher, err := GetBinaryPatcher(binaryTmpPath)
	if err != nil {
		return
	}
	return patcher.Apply()
}

// CanUpdate checks if the updater found a new version
func (u *Updater) CanUpdate() (bool, error) {
	latestVersion, err := u.Provider.GetLatestVersion()
	if err != nil {
		return false, err
	}
	if u.Version != latestVersion {
		return true, nil
	}
	return false, nil
}

// Update runs the updater
// It will update the current application if an update is found
func (u *Updater) Update() (err error) {
	canUpdate, err := u.CanUpdate()
	if err != nil || !canUpdate {
		return
	}
	if err = u.Provider.Open(); err != nil {
		return
	}
	defer u.Provider.Close()

	err = u.updateExecutable()
	return
}
