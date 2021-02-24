package provider_test

import (
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
