package provider_test

import (
	"path/filepath"
	"testing"

	provider "github.com/mouuff/go-rocket-update/provider"
)

func TestProviderZip(t *testing.T) {
	p, err := provider.NewZipProvider(filepath.Join("testdata", "Allum1.zip"))
	if err != nil {
		t.Error(err)
	}
	if err := p.Open(); err != nil {
		t.Error(err)
	}
	defer p.Close()

	err = ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Error(err)
	}
}
