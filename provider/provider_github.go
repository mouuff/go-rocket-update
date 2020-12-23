package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

// Github provider finds a zip file in the repository's releases to provide files
type Github struct {
	RepositoryURL string // Repository URL, example github.com/mouuff/go-rocket-update
	ZipName       string // Zip name (the zip you upload for a release on github), example: binaries.zip

	tmpDir      string // temporary directory this is used internally
	zipProvider *Zip   // provider used to unzip the downloaded zip
	zipPath     string // path to the downloaded zip (should be in tmpDir)
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

// getRepositoryInfo parses the github repository URL
func (c *Github) repositoryInfo() (*repositoryInfo, error) {
	re := regexp.MustCompile(`github\.com/(.*?)/(.*?)$`)
	submatches := re.FindAllStringSubmatch(c.RepositoryURL, 1)
	if len(submatches) < 1 {
		return nil, errors.New("Invalid github URL:" + c.RepositoryURL)
	}
	return &repositoryInfo{
		RepositoryOwner: submatches[0][1],
		RepositoryName:  submatches[0][2],
	}, nil
}

// getTagsURL get the tags URL for the github repository
func (c *Github) getTagsURL() (string, error) {
	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags",
		info.RepositoryOwner,
		info.RepositoryName,
	), nil
}

// getZipURL get the zip URL for the github repository
// If no tag is provided then the latest version is selected
func (c *Github) getZipURL(tag string) (string, error) {
	if len(tag) == 0 {
		// Get lastest version if no tag is provided
		var err error
		tag, err = c.GetLatestVersion()
		if err != nil {
			return "", err
		}
	}

	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		info.RepositoryOwner,
		info.RepositoryName,
		tag,
		c.ZipName,
	), nil
}

// getTags gets tags of the repository
func (c *Github) getTags() (tags []githubTag, err error) {
	tagsURL, err := c.getTagsURL()
	if err != nil {
		return
	}
	resp, err := http.Get(tagsURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return
	}
	return
}

// Open opens the provider
func (c *Github) Open() (err error) {
	zipURL, err := c.getZipURL("") // get zip url for lastest version
	if err != nil {
		return
	}
	resp, err := http.Get(zipURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	c.tmpDir, err = ioutil.TempDir("", "Github")
	if err != nil {
		return
	}

	c.zipPath = filepath.Join(c.tmpDir, c.ZipName)
	zipFile, err := os.Create(c.zipPath)
	if err != nil {
		return
	}
	_, err = io.Copy(zipFile, resp.Body)
	zipFile.Close()
	if err != nil {
		return
	}
	c.zipProvider = &Zip{Path: c.zipPath}
	return c.zipProvider.Open()
}

// Close closes the provider
func (c *Github) Close() error {
	if c.zipProvider != nil {
		c.zipProvider.Close()
		c.zipProvider = nil
	}

	if len(c.tmpDir) > 0 {
		os.RemoveAll(c.tmpDir)
		c.tmpDir = ""
		c.zipPath = ""
	}
	return nil
}

// GetLatestVersion gets the lastest version
func (c *Github) GetLatestVersion() (string, error) {
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
func (c *Github) Walk(walkFn WalkFunc) error {
	return c.zipProvider.Walk(walkFn)
}

// Retrieve file relative to "provider" to destination
func (c *Github) Retrieve(src string, dest string) error {
	return c.zipProvider.Retrieve(src, dest)
}
