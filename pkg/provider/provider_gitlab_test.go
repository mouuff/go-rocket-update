package provider_test

import (
	"log"
	"runtime"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGitlab(t *testing.T) {
	p := &provider.Gitlab{
		ProjectID: 24021648,
		ZipName:   "binaries_" + runtime.GOOS + ".zip",
	}

	log.Println(p.GetLatestVersion())

	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}
}
