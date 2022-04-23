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

	if err := p.Retrieve("x", "x"); err == nil {
		t.Fatal("Retrieve should return an error")
	}

	if err := p.Open(); err != nil {
		t.Fatal(err)
	}

	// Call to open twice should not create an error
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
	err = ProviderTestUnavailable(badProvider)
	if err != nil {
		t.Fatal(err)
	}
}
