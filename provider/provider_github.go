package provider

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
func (c *providerGithub) repositoryInfo() (*repositoryInfo, error) {
	re := regexp.MustCompile(`github\.com/(.*?)/(.*?)$`)
	submatches := re.FindAllStringSubmatch(c.repoURL, -1)
	if len(submatches) < 1 {
		return nil, errors.New("Invalid github URL:" + c.repoURL)
	}
	matches := submatches[0]
	if len(matches) != 3 {
		return nil, errors.New("Invalid github URL:" + c.repoURL)
	}
	return &repositoryInfo{
		RepositoryOwner: matches[1],
		RepositoryName:  matches[2],
	}, nil
}

// tagsURL get the tags URL for the github repository
func (c *providerGithub) tagsURL() (string, error) {
	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags",
		info.RepositoryOwner,
		info.RepositoryName,
	), nil
}

// zipURL get the zip URL
func (c *providerGithub) zipURL(tag string) (string, error) {
	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		info.RepositoryOwner,
		info.RepositoryName,
		tag,
		c.zipName,
	), nil
}

// getTags gets tags of the repository
func (c *providerGithub) getTags() (string, error) {
	tagsURL, err := c.tagsURL()
	if err != nil {
		return "", err
	}
	fmt.Println(tagsURL)
	response, err := http.Get(tagsURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return "", nil
}

// Open opens the provider
func (c *providerGithub) Open() error {
	_, err := c.getTags()
	return err
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
