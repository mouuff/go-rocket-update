package provider

import (
	"errors"
	"regexp"
)

type providerGithub struct {
	repoURL     string
	zipName     string
	zipProvider *providerZip
}

type repositoryInfo struct {
	RepositoryOwner string
	RepositoryName  string
}

// NewProviderGithub creates a new provider for local files
func NewProviderGithub(repoURL, zipName string) Provider {
	return &providerGithub{
		repoURL: repoURL,
		zipName: zipName,
	}
}

// getRepositoryInfo parses the github repository URL
func (c *providerGithub) RepositoryInfo() (*repositoryInfo, error) {
	re := regexp.MustCompile(`github\.com/(.*?)/(.*?)$`)
	matches := re.FindAllString(c.repoURL, -1)
	if len(matches) != 2 {
		return nil, errors.New("Invalid github URL:" + c.repoURL)
	}
	return &repositoryInfo{
		RepositoryOwner: matches[0],
		RepositoryName:  matches[1],
	}, nil
}

// Open opens the provider
func (c *providerGithub) Open() error {
	return nil
}

// Close closes the provider
func (c *providerGithub) Close() error {
	return nil
}

// GetLatestVersion gets the lastest version
func (c *providerGithub) GetLatestVersion() (string, error) {
	return "1.0", nil
}

// Walk walks all the files provided
func (c *providerGithub) Walk(walkFn WalkFunc) error {
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *providerGithub) Retrieve(src string, dest string) error {
	return nil
}
