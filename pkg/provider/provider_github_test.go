package provider_test

import (
	"runtime"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGithub(t *testing.T) {
	p := &provider.Github{
		RepositoryURL: "github.com/mouuff/go-rocket-update-example",
		ZipName:       "binaries_" + runtime.GOOS + ".zip",
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}
}
