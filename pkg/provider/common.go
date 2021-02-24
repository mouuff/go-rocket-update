package provider

import (
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// GetLatestVersionFromPath finds the latest version from a path
// This is used by provider zip and gzip
// Example: /example/binaries-v1.4.53.zip is going to match "v1.4.53"
func GetLatestVersionFromPath(path string) (string, error) {
	// TODO unittest
	versionRegex := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
	version := string(versionRegex.Find([]byte(filepath.Base(path))))
	if version == "" {
		return "", ErrProviderUnavaiable
	}
	return version, nil
}

// GlobNewestFile same as filepath.Glob but returns only one file with the latest modification time
func GlobNewestFile(pattern string) (string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	var newestTime time.Time
	var newestFile string
	for _, match := range matches {
		file, err := os.Stat(match)
		if err == nil {
			if newestFile == "" || newestTime.Before(file.ModTime()) {
				newestFile = match
				newestTime = file.ModTime()
			}
		}
	}
	if newestFile == "" {
		return "", ErrFileNotFound
	}
	return newestFile, nil
}
