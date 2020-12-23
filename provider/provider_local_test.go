package provider_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/easy-update/fileio"
	provider "github.com/mouuff/easy-update/provider"
)

func TestProviderLocal(t *testing.T) {
	p := &provider.Local{Path: filepath.Join("testdata", "Allum1")}
	if err := p.Open(); err != nil {
		t.Error(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Error(err)
	}

	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	destPath := filepath.Join(tmpDir, "test.txt")
	err = p.Retrieve(filepath.Join("subfolder", "testfile.txt"), destPath)
	if err != nil {
		t.Error(err)
	}
	equals, err := fileio.CompareFileChecksum(destPath, filepath.Join("testdata", "Allum1", "subfolder", "testfile.txt"))
	if err != nil {
		t.Error(err)
	}
	if equals == false {
		t.Error("Files should be equals")
	}
}
