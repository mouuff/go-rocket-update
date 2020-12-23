package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type providerGithub struct {
	repoURL     string
	zipName     string
	zipProvider *providerZip
}

// githubTag struct used to unmarshal response from github
// https://api.github.com/repos/ownerName/projectName/tags
type githubTag struct {
	Name string `json:"name"`
}

// repositoryInfo is used to get the name of the project and the owner name
// from this fields we are able to get other links (such as the release and tags link)
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
	submatches := re.FindAllStringSubmatch(c.repoURL, 1)
	if len(submatches) < 1 {
		return nil, errors.New("Invalid github URL:" + c.repoURL)
	}
	return &repositoryInfo{
		RepositoryOwner: submatches[0][1],
		RepositoryName:  submatches[0][2],
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
func (c *providerGithub) getTags() ([]githubTag, error) {
	var tags []githubTag
	tagsURL, err := c.tagsURL()
	if err != nil {
		return tags, err
	}
	response, err := http.Get(tagsURL)
	if err != nil {
		return tags, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&tags)
	if err != nil {
		return tags, err
	}
	return tags, nil
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
	tags, err := c.getTags()
	if err != nil {
		return "", err
	}
	if len(tags) < 1 {
		return "", errors.New("This github project has no tags")
	}
	return tags[0].Name, nil
}

// Walk walks all the files provided
func (c *providerGithub) Walk(walkFn WalkFunc) error {
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *providerGithub) Retrieve(src string, dest string) error {
	return nil
}
