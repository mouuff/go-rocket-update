package provider_test

import (
	"fmt"
	"runtime"
	"testing"

	provider "github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderGitlab(t *testing.T) {
	p := &provider.Gitlab{
		ProjectID:   24021648,
		ArchiveName: fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
	}

	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}

	badProvider := &provider.Gitlab{
		ProjectID:   424242424242424242,
		ArchiveName: fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
	}
	err = ProviderTestUnavailable(badProvider)
	if err != nil {
		t.Fatal(err)
	}
}
