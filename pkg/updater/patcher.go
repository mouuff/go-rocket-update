package updater

import (
	"io/ioutil"
	"os"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// ReplaceWith replaces file located at 'dest' with the one located at 'src'
// the 'dest' file is moved to 'backup' if 'backup' already exists it is replaced
// 'src' is leaved untouched
// Do not use this method for large files!

// VerifyFunc is used to verify that the
type VerifyFunc func(path string) error

// Patcher is used to replace file located at DestinationPath with the
// one located at SourcePath. DestinationPath is backed up at BackupPath.
//
type Patcher struct {
	DestinationPath string
	SourcePath      string
	BackupPath      string
	Mode            os.FileMode // Mode used to create the new file at DestinationPath
}

// Apply replaces the file located at DestinationPath with the one located at BackupPath and
// then new file is created at DestinationPath with the content of SourcePath
func (p *Patcher) Apply() error {
	content, err := ioutil.ReadFile(p.SourcePath)
	if err != nil {
		return err
	}
	if fileio.FileExists(p.BackupPath) {
		os.Remove(p.BackupPath)
	}
	err = os.Rename(p.DestinationPath, p.BackupPath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(p.DestinationPath, content, p.Mode)
	if err != nil {
		p.Rollback()
		return err
	}
	return nil
}

// Rollback replaces the file located at BackupPath with the one located at DestinationPath.
func (p *Patcher) Rollback() error {
	return os.Rename(p.BackupPath, p.DestinationPath)
}
