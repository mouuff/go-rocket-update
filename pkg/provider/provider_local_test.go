package provider_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/fileio"
	"github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderLocal(t *testing.T) {
	p := &provider.Local{Path: filepath.Join("testdata", "Allum1")}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	destPath := filepath.Join(tmpDir, "test.txt")
	err = p.Retrieve(filepath.Join("subfolder", "testfile.txt"), destPath)
	if err != nil {
		t.Fatal(err)
	}
	equals, err := fileio.CompareFiles(destPath, filepath.Join("testdata", "Allum1", "subfolder", "testfile.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if equals == false {
		t.Fatal("Files should be equals")
	}

	p = &provider.Local{Path: filepath.Join("testdata", "doesnotexists")}
	if err := p.Open(); err == nil {
		t.Error("Open() should return an error when path does not exists")
	}
	_, err = p.GetLatestVersion()
	if err == nil {
		t.Error("GetLatestVersion() should return an error when path does not exists")
	}

}
