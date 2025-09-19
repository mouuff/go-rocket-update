package updater

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mouuff/go-rocket-update/internal/fileio"
	"github.com/mouuff/go-rocket-update/pkg/provider"
)

// UpdateStatus represents the status after Updater{}.Update() was called
type UpdateStatus int

const (
	// Unknown update status (something went wrong)
	Unknown UpdateStatus = iota
	// UpToDate means the software is already up to date
	UpToDate
	// Updated means the software have been updated
	Updated
)

// PostUpdateFunc is called after a successful update
// On error you can/should call u.Rollback() yourself
type PostUpdateFunc func(u *Updater) (UpdateStatus, error)

// Updater struct
type Updater struct {
	Provider           provider.Provider
	ExecutableName     string         // Name of the executable
	Version            string         // The current version of your program
	OverrideExecutable string         // (optional) Overrides the path of the executable
	PostUpdateFunc     PostUpdateFunc // (optional) Set a function that will be called after an update (see type documentation)
	latestVersion      string         // cache for the latest version
}

// getExecutablePatcher gets the executable patcher
// executableCandidate can be empty if you only plan to rollback
func (u *Updater) getExecutablePatcher(executableCandidatePath string) (*fileio.Patcher, error) {
	executable, err := u.GetExecutable()
	if err != nil {
		return nil, err
	}
	return &fileio.Patcher{

		SourcePath:      executableCandidatePath,
		DestinationPath: executable,
		BackupPath:      executable + ".old",
		Mode:            0755,
	}, nil
}

// findExecutableRemotePath finds the remote executable path using the provider
func (u *Updater) findExecutableRemotePath() (string, error) {
	executableRemotePath := ""
	err := u.Provider.Walk(func(info *provider.FileInfo) error {
		if info.Mode.IsRegular() && strings.Contains(filepath.Base(info.Path), u.ExecutableName) {
			executableRemotePath = info.Path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return executableRemotePath, nil
}

// updateExecutable updates the current executable with the new one
func (u *Updater) updateExecutable() (err error) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)
	executableRemotePath, err := u.findExecutableRemotePath()
	if err != nil {
		return
	}
	executableCanditatePath := filepath.Join(tmpDir, filepath.Base(executableRemotePath))
	err = u.Provider.Retrieve(executableRemotePath, executableCanditatePath)
	if err != nil {
		return
	}

	patcher, err := u.getExecutablePatcher(executableCanditatePath)
	if err != nil {
		return
	}
	return patcher.Apply() // on failure it will automatically rollback already
}

// GetExecutable gets the executable path that will be used to for the update process
// same as fileio.GetExecutable() but this one takes into account the variable OverrideExecutablePath
func (u *Updater) GetExecutable() (string, error) {
	if u.OverrideExecutable == "" {
		return fileio.GetExecutable()
	}
	return u.OverrideExecutable, nil
}

// GetLatestVersion gets the latest version (same as provider.GetLatestVersion but keeps the version in cache)
func (u *Updater) GetLatestVersion() (string, error) {
	if u.latestVersion != "" {
		return u.latestVersion, nil
	}
	var err error
	u.latestVersion, err = u.Provider.GetLatestVersion()
	if err != nil {
		u.latestVersion = ""
		return u.latestVersion, err
	}
	return u.latestVersion, nil
}

// CanUpdate checks if the updater found a new version
func (u *Updater) CanUpdate() (bool, error) {
	latestVersion, err := u.GetLatestVersion()
	if err != nil {
		return false, err
	}
	if u.Version != latestVersion {
		return true, nil
	}
	return false, nil
}

// Update runs the updater
// It will update the current application if an update is found.
// If the application is up-to-date then no error will be returned.
// YOU DON'T NEED TO call Rollback() yourself!
// UNLESS you want to rollback after a successful update
func (u *Updater) Update() (status UpdateStatus, err error) {
	status = Unknown
	canUpdate, err := u.CanUpdate()
	if err != nil {
		return
	}
	if !canUpdate {
		status = UpToDate
		return
	}
	if err = u.Provider.Open(); err != nil {
		return
	}
	defer u.Provider.Close()

	err = u.updateExecutable()
	if err == nil {
		status = Updated
	}
	// WARNING: any code after that should also call Rollback() on failure
	if status == Updated && u.PostUpdateFunc != nil {
		status, err = u.PostUpdateFunc(u) // That's up to the PostUpdateFunc to rollback or not here
	}
	return
}

// Rollback rollbacks to the previous version
// Use this if you know what you are doing
func (u *Updater) Rollback() (err error) {
	executablePatcher, err := u.getExecutablePatcher("")
	if err != nil {
		return err
	}
	return executablePatcher.Rollback()
}

// CleanUp removes the old executable after updating. (.old files)
// Warning: This also prevents any rollback
func (u *Updater) CleanUp() (err error) {
	executablePatcher, err := u.getExecutablePatcher("")
	if err != nil {
		return err
	}

	return executablePatcher.CleanUp()
}
