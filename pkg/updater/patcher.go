package updater

import (
	"io/ioutil"
	"os"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// VerifyFunc is used to verify that the file is valid
type VerifyFunc func(path string) error

// Patcher is used to replace file located at DestinationPath with the
// one located at SourcePath. DestinationPath is backed up at BackupPath.
type Patcher struct {
	DestinationPath string
	SourcePath      string
	BackupPath      string
	Mode            os.FileMode // Mode used to create the new file at DestinationPath
	Verify          VerifyFunc  // Verify (optionnal) is a function to verify that the installation is good
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
	if p.Verify != nil {
		if err := p.Verify(p.DestinationPath); err != nil {
			return p.Rollback()
		}
	}
	return nil
}

// Rollback replaces the file located at BackupPath with the one located at DestinationPath.
func (p *Patcher) Rollback() error {
	return os.Rename(p.BackupPath, p.DestinationPath)
}

// CleanUp cleans up backup file
func (p *Patcher) CleanUp() error {
	if fileio.FileExists(p.BackupPath) {
		return os.Remove(p.BackupPath)
	}
	return nil
}
