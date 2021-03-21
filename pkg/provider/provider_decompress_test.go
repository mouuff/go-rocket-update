package provider_test

import (
	"path/filepath"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderDecompressZip(t *testing.T) {
	p, err := provider.Decompress(filepath.Join("testdata", "Allum1-v1.0.0.zip"))
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err = ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProviderDecompressGzip(t *testing.T) {
	p, err := provider.Decompress(filepath.Join("testdata", "Allum1-v1.0.0.tar.gz"))
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err = ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProviderDecompressUnknown(t *testing.T) {
	_, err := provider.Decompress(filepath.Join("testdata", "Allum1-v1.0.0.x.x"))
	if err == nil {
		t.Error("Should return an error if format to decompress is unknown")
	}
}
