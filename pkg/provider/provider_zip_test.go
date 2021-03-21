package provider_test

import (
	"path/filepath"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderZip(t *testing.T) {
	p := &provider.Zip{
		Path: filepath.Join("testdata", "Allum1-v1.0.0.zip"),
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}

	badProvider := &provider.Zip{
		Path: filepath.Join("testdata", "unknownpath.zip"),
	}
	err = ProviderTestUnavaiable(badProvider)
	if err != nil {
		t.Fatal(err)
	}

}
