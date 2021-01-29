package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Gitlab provider finds a zip file in the repository's releases to provide files
type Gitlab struct {
	RepositoryURL string // Repository URL, example gitlab.com/mouuff/go-rocket-update
	ZipName       string // Zip name (the zip you upload for a release on gitlab), example: binaries.zip

	tmpDir      string // temporary directory this is used internally
	zipProvider *Zip   // provider used to unzip the downloaded zip
	zipPath     string // path to the downloaded zip (should be in tmpDir)
}

// gitlabTag struct used to unmarshal response from gitlab
// https://api.gitlab.com/repos/ownerName/projectName/tags
type gitlabTag struct {
	Name string `json:"name"`
}

// repositoryInfo is used to get the name of the project and the owner name
// from this fields we are able to get other links (such as the release and tags link)
type repositoryInfo struct {
	RepositoryOwner string
	RepositoryName  string
}

// getRepositoryInfo parses the gitlab repository URL
func (c *Gitlab) repositoryInfo() (*repositoryInfo, error) {
	re := regexp.MustCompile(`gitlab\.com/(.*?)/(.*?)$`)
	submatches := re.FindAllStringSubmatch(c.RepositoryURL, 1)
	if len(submatches) < 1 {
		return nil, errors.New("Invalid gitlab URL:" + c.RepositoryURL)
	}
	return &repositoryInfo{
		RepositoryOwner: submatches[0][1],
		RepositoryName:  submatches[0][2],
	}, nil
}

// getTagsURL get the tags URL for the gitlab repository
func (c *Gitlab) getTagsURL() (string, error) {
	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.gitlab.com/repos/%s/%s/tags",
		info.RepositoryOwner,
		info.RepositoryName,
	), nil
}

// getZipURL get the zip URL for the gitlab repository
// If no tag is provided then the latest version is selected
func (c *Gitlab) getZipURL(tag string) (string, error) {
	if len(tag) == 0 {
		// Get latest version if no tag is provided
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
	return fmt.Sprintf("https://gitlab.com/%s/%s/releases/download/%s/%s",
		info.RepositoryOwner,
		info.RepositoryName,
		tag,
		c.ZipName,
	), nil
}

// getTags gets tags of the repository
func (c *Gitlab) getTags() (tags []gitlabTag, err error) {
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
func (c *Gitlab) Open() (err error) {
	zipURL, err := c.getZipURL("") // get zip url for latest version
	if err != nil {
		return
	}
	resp, err := http.Get(zipURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	c.tmpDir, err = fileio.TempDir()
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
func (c *Gitlab) Close() error {
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

// GetLatestVersion gets the latest version
func (c *Gitlab) GetLatestVersion() (string, error) {
	tags, err := c.getTags()
	if err != nil {
		return "", err
	}
	if len(tags) < 1 {
		return "", errors.New("This gitlab project has no tags")
	}
	return tags[0].Name, nil
}

// Walk walks all the files provided
func (c *Gitlab) Walk(walkFn WalkFunc) error {
	return c.zipProvider.Walk(walkFn)
}

// Retrieve file relative to "provider" to destination
func (c *Gitlab) Retrieve(src string, dest string) error {
	return c.zipProvider.Retrieve(src, dest)
}
