package provider_test

import (
	"path/filepath"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestGetLatestVersionFromPath(t *testing.T) {
	version, err := provider.GetLatestVersionFromPath("binaries-v1.2.3.zip")
	if err != nil {
		t.Fatal(err)
	}
	if version != "v1.2.3" {
		t.Error("version != 'v1.2.3'")
	}

	version, err = provider.GetLatestVersionFromPath("binaries-v.zip")
	if err == nil {
		t.Error("Should return an error")
	}
}

func TestGlobNewestFile(t *testing.T) {
	match, err := provider.GlobNewestFile(filepath.Join("testdata", "Allum1-v*.tar.gz"))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(match) != "Allum1-v1.1.0.tar.gz" {
		t.Error("Expected Allum1-v1.1.0.tar.gz")
	}

	match, err = provider.GlobNewestFile(filepath.Join("testdata", "doesntexists"))
	if err == nil {
		t.Error("Should return an error if file doesn't exists")
	}
}
