package provider_test

import (
	"fmt"
	"runtime"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGithub(t *testing.T) {
	p := &provider.Github{
		RepositoryURL: "github.com/mouuff/go-rocket-update-example",
		ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}

	badProvider := &provider.Github{
		RepositoryURL: "github.com/XXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXXXXX",
		ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
	}
	err = ProviderTestUnavailable(badProvider)
	if err != nil {
		t.Fatal(err)
	}
}
