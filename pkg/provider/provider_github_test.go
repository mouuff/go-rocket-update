package provider_test

import (
	"os"
	"runtime"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGithub(t *testing.T) {
	if os.Getenv("CI") != "" {
		// We skip this test because using travis the request to github often fails with a "HTTP 429 Too Many Requests" error
		t.Skip("Skipping testing in CI environment")
	}
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
