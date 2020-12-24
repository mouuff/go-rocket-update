package resolver

import "github.com/mouuff/go-rocket-update/provider"

// A Resolver must define what file to retrieve in the provider
// This is used to separate the file lookup logic from the updater
type Resolver interface {
	SetProvider(provider provider.Provider)
	RetrieveBinary(dest string) error // Download the binary from the provider
}
