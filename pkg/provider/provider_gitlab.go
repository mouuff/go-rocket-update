package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Gitlab provider finds a zip file in the repository's releases to provide files
type Gitlab struct {
	RepositoryURL string // Repository URL, example gitlab.com/mouuff/go-rocket-update
	ProjectID     int
	ZipName       string // Zip name (the zip you upload for a release on gitlab), example: binaries.zip

	tmpDir      string // temporary directory this is used internally
	zipProvider *Zip   // provider used to unzip the downloaded zip
	zipPath     string // path to the downloaded zip (should be in tmpDir)
}

// gitlabRelease struct used to unmarshal response from gitlab
// https://gitlab.com/api/v4/projects/24021648/releases
type gitlabRelease struct {
	Name string `json:"tag_name"`
}

// getReleasesURL get the releases URL for the gitlab repository
func (c *Gitlab) getReleasesURL() (string, error) {
	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/releases",
		c.ProjectID,
	), nil
}

// getZipURL get the zip URL for the gitlab repository
// If no tag is provided then the latest version is selected
func (c *Gitlab) getZipURL(tag string) (string, error) {
	/*
		if len(tag) == 0 {
			// Get latest version if no tag is provided
			var err error
			tag, err = c.GetLatestVersion()
			if err != nil {
				return "", err
			}
		}
		return fmt.Sprintf("https://gitlab.com/%s/%s/releases/download/%s/%s",
			info.RepositoryOwner,
			info.RepositoryName,
			tag,
			c.ZipName,
		), nil
	*/
	return "", nil
}

// getReleases gets tags of the repository
func (c *Gitlab) getReleases() (tags []gitlabRelease, err error) {
	releasesURL, err := c.getReleasesURL()
	if err != nil {
		return
	}
	resp, err := http.Get(releasesURL)
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
	tags, err := c.getReleases()
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
