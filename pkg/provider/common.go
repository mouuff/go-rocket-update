package provider

import (
	"path/filepath"
	"regexp"
)

// getLatestVersionFromPath finds the latest version from a path
// This is used by provider zip and gzip
// Example: /example/binaries-v1.4.53.zip is going to match "v1.4.53"
func getLatestVersionFromPath(path string) (string, error) {
	// TODO unittest
	versionRegex := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
	version := string(versionRegex.Find([]byte(filepath.Base(path))))
	if version == "" {
		return "", ErrProviderUnavaiable
	}
	return version, nil
}
