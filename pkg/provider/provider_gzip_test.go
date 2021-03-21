package provider_test

import (
	"path/filepath"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGzip(t *testing.T) {
	p := &provider.Gzip{
		Path: filepath.Join("testdata", "Allum1-v1.0.0.tar.gz"),
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}
	badProvider := &provider.Gzip{
		Path: filepath.Join("testdata", "doesnotexist.tar.gz"),
	}
	err = ProviderTestUnavaiable(badProvider)
	if err != nil {
		t.Fatal(err)
	}
}
