package updater_test

import (
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/fileio"
	"github.com/mouuff/go-rocket-update/pkg/updater"
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

	sourcePath := filepath.Join(tmpDir, "file.jpeg")
	destinationPath := filepath.Join(tmpDir, "binary")
	backupPath := sourcePath + ".old"

	err = fileio.CopyFile(
		filepath.Join("testdata", "binary"),
		destinationPath)
	if err != nil {
		t.Fatal(err)
	}
	err = fileio.CopyFile(
		filepath.Join("testdata", "file.jpeg"),
		sourcePath)
	if err != nil {
		t.Fatal(err)
	}

	patcher := &updater.Patcher{
		DestinationPath: destinationPath,
		SourcePath:      sourcePath,
		BackupPath:      backupPath,
		Mode:            0755,
		Verify:          nil, // TODO
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
}
