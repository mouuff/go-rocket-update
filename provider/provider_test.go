package providers_test

import (
	"path"
	"testing"

	provider "github.com/mouuff/easy-update/provider"
)

func testProvider(p provider.Provider) error {
	return nil
}

func TestProviderLocal(t *testing.T) {
	p := provider.NewProviderLocal(path.Join("testdata", "Allum1"))
}
