package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Gitlab provider finds a zip file in the repository's releases to provide files
type Gitlab struct {
	ProjectID int
	ZipName   string // Zip name (the zip you upload for a release on gitlab), example: binaries.zip

	tmpDir      string // temporary directory this is used internally
	zipProvider *Zip   // provider used to unzip the downloaded zip
	zipPath     string // path to the downloaded zip (should be in tmpDir)
}

// gitlabRelease struct used to unmarshal response from gitlab
// https://gitlab.com/api/v4/projects/24021648/releases
type gitlabRelease struct {
	TagName string              `json:"tag_name"`
	Links   []gitlabReleaseLink `json:"links"`
}

type gitlabReleaseLink struct {
	Name      string `json:"name"`
	DirectURL string `json:"direct_asset_url"`
}

// getReleasesURL get the releases URL for the gitlab repository
func (c *Gitlab) getReleasesURL() (string, error) {
	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/releases",
		c.ProjectID,
	), nil
}

// getZipURL get the zip URL for the gitlab repository
// the latest version is selected
func (c *Gitlab) getZipURL() (string, error) {
	release, err := c.getLatestRelease()
	if err != nil {
		return "", err
	}
	log.Println(release)
	for _, link := range release.Links {
		log.Println(link)
		if strings.HasSuffix(link.Name, c.ZipName) {
			return link.DirectURL, nil
		}
	}
	return "", errors.New("Link not found for name: " + c.ZipName)
}

// getReleases gets tags of the repository
func (c *Gitlab) getReleases() (releases []gitlabRelease, err error) {
	releasesURL, err := c.getReleasesURL()
	if err != nil {
		return
	}
	resp, err := http.Get(releasesURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return
	}
	return
}

// getReleases gets tags of the repository
func (c *Gitlab) getLatestRelease() (*gitlabRelease, error) {
	releases, err := c.getReleases()
	if err != nil {
		return nil, err
	}
	if len(releases) < 1 {
		return nil, errors.New("This gitlab project has no releases")
	}
	return &releases[0], nil
}

// Open opens the provider
func (c *Gitlab) Open() (err error) {
	zipURL, err := c.getZipURL() // get zip url for latest version
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
	release, err := c.getLatestRelease()
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

// Walk walks all the files provided
func (c *Gitlab) Walk(walkFn WalkFunc) error {
	return c.zipProvider.Walk(walkFn)
}

// Retrieve file relative to "provider" to destination
func (c *Gitlab) Retrieve(src string, dest string) error {
	return c.zipProvider.Retrieve(src, dest)
}
