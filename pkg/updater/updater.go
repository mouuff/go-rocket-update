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
	Provider           provider.Provider
	BinaryName         string
	Version            string
	OverrideExecutable string // (optionnal) Overrides the path of the executable
}

// GetExecutable gets the executable path that will be used to for the update process
// same as fileio.GetExecutable() but this one takes into account the variable OverrideExecutablePath
func (u *Updater) GetExecutable() (string, error) {
	if u.OverrideExecutable == "" {
		return fileio.GetExecutable()
	}
	return u.OverrideExecutable, nil
}

// getBinaryPatcher gets the binary patcher
// binaryCandidate can be empty if you only plan to rollback
func (u *Updater) getBinaryPatcher(binaryCandidatePath string) (*fileio.Patcher, error) {
	executable, err := u.GetExecutable()
	if err != nil {
		return nil, err
	}
	return &fileio.Patcher{

		SourcePath:      binaryCandidatePath,
		DestinationPath: executable,
		BackupPath:      executable + ".old",
		Mode:            0755,
	}, nil
}

// getBinaryName gets the name used to find the right binary
func (u *Updater) getBinaryName() string {
	return u.BinaryName + "_" + runtime.GOOS + "_" + runtime.GOARCH
}

// findBinaryRemotePath finds the right binary using the provider
func (u *Updater) findBinaryRemotePath() (string, error) {
	binaryRemotePath := ""
	err := u.Provider.Walk(func(info *provider.FileInfo) error {
		if info.Mode.IsRegular() && strings.Contains(filepath.Base(info.Path), u.getBinaryName()) {
			binaryRemotePath = info.Path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return binaryRemotePath, nil
}

// updateExecutable updates the current executable with the new one
func (u *Updater) updateExecutable() (err error) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)
	binaryRemotePath, err := u.findBinaryRemotePath()
	if err != nil {
		return
	}
	binaryCanditatePath := filepath.Join(tmpDir, filepath.Base(binaryRemotePath))
	err = u.Provider.Retrieve(binaryRemotePath, binaryCanditatePath)
	if err != nil {
		return
	}

	patcher, err := u.getBinaryPatcher(binaryCanditatePath)
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

// Rollback rollbacks to the previous version
// Use this if you know what you are doing
func (u *Updater) Rollback() (err error) {
	binaryPatcher, err := u.getBinaryPatcher("")
	if err != nil {
		return err
	}
	return binaryPatcher.Rollback()
}
