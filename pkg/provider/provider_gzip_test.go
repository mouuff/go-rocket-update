package provider_test

import (
	"fmt"
	"path/filepath"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGzip(t *testing.T) {
	p := &provider.Gzip{
		Path: filepath.Join("testdata", "Allum1.tar.gz"),
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	err := p.Walk(func(info *provider.FileInfo) error {
		fmt.Println(info.Path)
		fmt.Println(info.Mode.IsDir())
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	/*
		err := ProviderTestWalkAndRetrieve(p)
		if err != nil {
			t.Fatal(err)
		}
	*/
}
