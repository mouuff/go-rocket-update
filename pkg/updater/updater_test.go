package updater_test

import (
	"encoding/hex"
	"fmt"
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

	executable := filepath.Join(tmpDir, "executable")
	err = fileio.CopyFile(filepath.Join("testdata", "testBinary"), executable)
	if err != nil {
		t.Fatal(err)
	}
	solutionDir := filepath.Join("testdata", "testSolution")
	u := &updater.Updater{
		Provider:           &provider.Local{Path: solutionDir},
		ExecutableName:     "test",
		Version:            "v1.0",
		OverrideExecutable: executable,
	}

	canUpdate, err := u.CanUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if canUpdate {
		t.Error("Should not be able to update with same version")
	}
	updateStatus, err := u.Update()
	if err != nil {
		t.Error(err)
	}
	if updateStatus != updater.UpToDate {
		t.Error("updateStatus should be updater.UpToDate if trying to update with the same version")
	}

	u.Version = "v0.1"

	canUpdate, err = u.CanUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if !canUpdate {
		t.Error("Should be able to update with different version")
	}
	executableChecksum, err := crypto.ChecksumFileSHA256(executable)
	if err != nil {
		t.Fatal(err)
	}
	updateStatus, err = u.Update()
	if err != nil {
		t.Fatal(err)
	}
	if updateStatus != updater.Updated {
		t.Error("updateStatus != updater.Updated")
	}
	updatedExecutableChecksum, err := crypto.ChecksumFileSHA256(executable)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(executableChecksum) == hex.EncodeToString(updatedExecutableChecksum) {
		t.Error("executableChecksum == updatedExecutableChecksum")
	}
	err = u.Rollback()
	if err != nil {
		t.Fatal(err)
	}
	rollbackedExecutableChecksum, err := crypto.ChecksumFileSHA256(executable)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(executableChecksum) != hex.EncodeToString(rollbackedExecutableChecksum) {
		t.Error("executableChecksum != rollbackedExecutableChecksum")
	}

	// Test with unreachable provider
	u = &updater.Updater{
		Provider:           &provider.Local{Path: solutionDir + "doesnotexit"},
		ExecutableName:     "test",
		Version:            "v1.0",
		OverrideExecutable: executable,
	}

	updateStatus, err = u.Update()
	if err == nil {
		t.Error("u.Update() should return an error when provider is not reachable")
	}
	if updateStatus == updater.Updated {
		t.Error("updateStatus == updater.Updated while provider is not reachable")
	}
	if updateStatus == updater.UpToDate {
		t.Error("updateStatus == updater.Updated while provider is not reachable")
	}
}

type TestHook struct {
	WillFail bool             // If set to true then PostUpdate will return an error
	Updater  *updater.Updater // This will be set when PostUpdate is called
}

func (hook *TestHook) PostUpdate(u *updater.Updater) error {
	hook.Updater = u
	if hook.WillFail {
		return fmt.Errorf("hook.WillFail")
	} else {
		return nil
	}
}

func TestUpdaterHooks(t *testing.T) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	executable := filepath.Join(tmpDir, "executable")
	err = fileio.CopyFile(filepath.Join("testdata", "testBinary"), executable)
	if err != nil {
		t.Fatal(err)
	}
	solutionDir := filepath.Join("testdata", "testSolution")
	hook := &TestHook{WillFail: true}
	u := &updater.Updater{
		Provider:           &provider.Local{Path: solutionDir},
		ExecutableName:     "test",
		Version:            "v0.0",
		OverrideExecutable: executable,
		Hook:               hook,
	}

	status, err := u.Update()
	if err == nil {
		t.Error("If hook fails, then the update should fail")
	}
	if status == updater.Updated {
		t.Error("If hook fails, then the status should not be set to Updated")
	}
	if hook.Updater != u {
		t.Error("hook.Updater != u")
	}
	hook.Updater = nil

	hook.WillFail = false
	_, err = u.Update()
	if err != nil {
		t.Error(err)
	}
	if hook.Updater != u {
		t.Error("hook.Updater != u")
	}
}
