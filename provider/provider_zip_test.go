package provider_test

import (
	"path/filepath"
	"testing"

	provider "github.com/mouuff/easy-update/provider"
)

func TestProviderZip(t *testing.T) {
	p := &provider.Zip{Path: filepath.Join("testdata", "Allum1.zip")}
	if err := p.Open(); err != nil {
		t.Error(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Error(err)
	}
}
