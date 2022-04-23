package fileio_test

import (
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func assertFilesEquals(t *testing.T, fileA, fileB string) {
	equals, err := fileio.CompareFiles(fileA, fileB)
	if err != nil {
		t.Fatal(err)
	}
	if !equals {
		t.Fatal("Should be equals")
	}
}

func assertFilesNotEquals(t *testing.T, fileA, fileB string) {
	equals, err := fileio.CompareFiles(fileA, fileB)
	if err != nil {
		t.Fatal(err)
	}
	if equals {
		t.Fatal("Should not be equals")
	}
}

func TestPatcher(t *testing.T) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	unknownPath := filepath.Join("pathdoesnotexists", "pathdoesnotexists")
	originalSourcePath := filepath.Join("testdata", "binary")
	originalDestinationFile := filepath.Join("testdata", "file.jpeg")

	sourcePath := filepath.Join(tmpDir, "file.jpeg")
	destinationPath := filepath.Join(tmpDir, "binary")
	backupPath := sourcePath + ".old"

	err = fileio.CopyFile(
		originalDestinationFile,
		destinationPath)
	if err != nil {
		t.Fatal(err)
	}
	err = fileio.CopyFile(
		originalSourcePath,
		sourcePath)
	if err != nil {
		t.Fatal(err)
	}

	patcher := &fileio.Patcher{
		DestinationPath: destinationPath,
		SourcePath:      sourcePath,
		BackupPath:      backupPath,
		Mode:            0755,
	}
	assertFilesNotEquals(t, sourcePath, destinationPath)

	err = patcher.Apply()
	if err != nil {
		t.Fatal(err)
	}

	if !fileio.FileExists(backupPath) {
		t.Fatal("Backup file should be created")
	}

	assertFilesEquals(t, sourcePath, destinationPath)

	err = patcher.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	assertFilesNotEquals(t, sourcePath, destinationPath)
	assertFilesEquals(t, originalSourcePath, sourcePath)

	err = patcher.CleanUp()
	if err != nil {
		t.Fatal(err)
	}
	if fileio.FileExists(backupPath) {
		t.Fatal("Backup file should be cleaned")
	}
	err = patcher.CleanUp() // Cleaning up a second time shouldn't cause problem
	if err != nil {
		t.Fatal(err)
	}

	// Test apply then clean up without rollback
	err = patcher.Apply()
	if err != nil {
		t.Fatal(err)
	}
	err = patcher.CleanUp()
	if err != nil {
		t.Fatal(err)
	}
	if fileio.FileExists(backupPath) {
		t.Fatal("Backup file should be cleaned")
	}

	patcher.DestinationPath = "wrongpath"

	err = patcher.Apply()
	if err == nil {
		t.Error("Should return an error")
	}
	patcher.DestinationPath = destinationPath
	patcher.SourcePath = "wrongpath"
	if err == nil {
		t.Error("Should return an error")
	}

	patcher.DestinationPath = unknownPath
	err = patcher.Apply()
	if err == nil {
		t.Error("Apply() should not work with unknown DestinationPath")
	}
	patcher.DestinationPath = destinationPath
	patcher.SourcePath = unknownPath
	err = patcher.Apply()
	if err == nil {
		t.Error("Apply() should not work with unknown SourcePath")
	}
	patcher.SourcePath = sourcePath
	patcher.BackupPath = unknownPath
	err = patcher.Apply()
	if err == nil {
		t.Error("Apply() should not work with unknown BackupPath")
	}
}
