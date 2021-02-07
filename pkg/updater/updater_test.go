package updater_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
	provider "github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

func TestUpdater(t *testing.T) {

	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	originalExecutable := filepath.Join(tmpDir, "executable")
	err = fileio.CopyFile(filepath.Join("testdata", "testBinary"), originalExecutable)
	if err != nil {
		t.Fatal(err)
	}
	solutionDir := filepath.Join("testdata", "testSolution")
	u := &updater.Updater{
		Provider:           &provider.Local{Path: solutionDir},
		ExecutableName:     "test",
		Version:            "v1.0",
		OverrideExecutable: originalExecutable,
	}

	canUpdate, err := u.CanUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if canUpdate {
		t.Error("Should not be able to update with same version")
	}
	u.Version = "v0.1"

	canUpdate, err = u.CanUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if !canUpdate {
		t.Error("Should be able to update with different version")
	}
	originalExecutableChecksum, err := crypto.ChecksumFileSHA256(originalExecutable)
	if err != nil {
		t.Fatal(err)
	}
	if err := u.Update(); err != nil {
		t.Fatal(err)
	}
	updatedExecutableChecksum, err := crypto.ChecksumFileSHA256(originalExecutable)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(originalExecutableChecksum) == hex.EncodeToString(updatedExecutableChecksum) {
		t.Error("originalExecutableChecksum == updatedExecutableChecksum")
	}
	err = u.Rollback()
	if err != nil {
		t.Fatal(err)
	}
	rollbackedExecutableChecksum, err := crypto.ChecksumFileSHA256(originalExecutable)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(originalExecutableChecksum) != hex.EncodeToString(rollbackedExecutableChecksum) {
		t.Error("originalExecutableChecksum != rollbackedExecutableChecksum")
	}
}
