package updater_test

import (
	"os"
	"path/filepath"
	"testing"

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

	mockExecutable := filepath.Join(tmpDir, "executable")
	err = fileio.CopyFile(filepath.Join("testdata", "testBinary"), mockExecutable)
	if err != nil {
		t.Fatal(err)
	}
	solutionDir := filepath.Join("testdata", "testSolution")
	u := &updater.Updater{
		Provider:           &provider.Local{Path: solutionDir},
		ExecutableName:     "test",
		Version:            "v1.0",
		OverrideExecutable: mockExecutable,
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
	if err := u.Update(); err != nil {
		t.Fatal(err)
	}
}
