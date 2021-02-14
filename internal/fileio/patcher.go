package fileio

import (
	"io/ioutil"
	"os"
)

// Patcher is used to replace file located at DestinationPath with the
// one located at SourcePath. DestinationPath is backed up at BackupPath.
type Patcher struct {
	SourcePath      string
	DestinationPath string
	BackupPath      string      // Stores the backup of DestinationPath before replacement (must be on the same filesytem as DestinationPath)
	Mode            os.FileMode // Mode used to create the new file at DestinationPath
}

// Apply replaces the file located at DestinationPath with the one located at BackupPath and
// then new file is created at DestinationPath with the content of SourcePath
func (p *Patcher) Apply() error {
	content, err := ioutil.ReadFile(p.SourcePath)
	if err != nil {
		return err
	}
	_ = p.CleanUp() // We dont check error on purpose
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

// CleanUp cleans up backup file
func (p *Patcher) CleanUp() error {
	if FileExists(p.BackupPath) {
		return os.Remove(p.BackupPath)
	}
	return nil
}
