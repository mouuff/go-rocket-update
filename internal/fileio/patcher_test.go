package fileio_test

import (
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func AssertFilesEquals(t *testing.T, fileA, fileB string) {
	equals, err := fileio.CompareFiles(fileA, fileB)
	if err != nil {
		t.Fatal(err)
	}
	if !equals {
		t.Fatal("Should be equals")
	}
}

func AssertFilesNotEquals(t *testing.T, fileA, fileB string) {
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
	AssertFilesNotEquals(t, sourcePath, destinationPath)

	err = patcher.Apply()
	if err != nil {
		t.Fatal(err)
	}

	AssertFilesEquals(t, sourcePath, destinationPath)

	err = patcher.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	AssertFilesNotEquals(t, sourcePath, destinationPath)
	AssertFilesEquals(t, originalSourcePath, sourcePath)

	err = patcher.CleanUp()
	if err != nil {
		t.Fatal(err)
	}
	if fileio.FileExists(backupPath) {
		t.Fatal("Backup file should be cleaned")
	}
}
