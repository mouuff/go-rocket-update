package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Gitlab provider finds a archive file in the repository's releases to provide files
type Gitlab struct {
	ProjectID   int
	ArchiveName string // ArchiveName (the archive you upload for a release on gitlab), example: binaries.zip
	ApiURI      string // ApiURI (in case you're using a private gitlab server), example: gitlab.mydomain.tld/api/v4/projects/%d/releases to use gitlab.com let it blank

	tmpDir             string   // temporary directory this is used internally
	decompressProvider Provider // provider used to decompress the downloaded archive
	decompressPath     string   // path to the downloaded archive (should be in tmpDir)
}

// gitlabRelease struct used to unmarshal response from gitlab
// https://gitlab.com/api/v4/projects/24021648/releases
type gitlabRelease struct {
	TagName string               `json:"tag_name"`
	Assets  *gitlabReleaseAssets `json:"assets"`
}

type gitlabReleaseAssets struct {
	Links []gitlabReleaseLink `json:"links"`
}

type gitlabReleaseLink struct {
	Name      string `json:"name"`
	DirectURL string `json:"direct_asset_url"`
}

// getReleasesURL get the releases URL for the gitlab repository
func (c *Gitlab) getReleasesURL() (string, error) {
	if c.ApiURI != "" {
		return fmt.Sprintf(c.ApiURI,
			c.ProjectID,
		), nil
	}

	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/releases",
		c.ProjectID,
	), nil
}

// getArchiveURL get the archive URL for the gitlab repository
// the latest version is selected
func (c *Gitlab) getArchiveURL() (string, error) {
	release, err := c.getLatestRelease()
	if err != nil {
		return "", err
	}
	for _, link := range release.Assets.Links {
		if strings.HasSuffix(link.Name, c.ArchiveName) {
			return link.DirectURL, nil
		}
	}
	return "", errors.New("Link not found for name: " + c.ArchiveName)
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
	archiveURL, err := c.getArchiveURL() // get archive url for latest version
	if err != nil {
		return
	}
	resp, err := http.Get(archiveURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	c.tmpDir, err = fileio.TempDir()
	if err != nil {
		return
	}

	c.decompressPath = filepath.Join(c.tmpDir, c.ArchiveName)
	archiveFile, err := os.Create(c.decompressPath)
	if err != nil {
		return
	}
	_, err = io.Copy(archiveFile, resp.Body)
	archiveFile.Close()
	if err != nil {
		return
	}
	c.decompressProvider, err = Decompress(c.decompressPath)
	return c.decompressProvider.Open()
}

// Close closes the provider
func (c *Gitlab) Close() error {
	if c.decompressProvider != nil {
		c.decompressProvider.Close()
		c.decompressProvider = nil
	}

	if len(c.tmpDir) > 0 {
		os.RemoveAll(c.tmpDir)
		c.tmpDir = ""
		c.decompressPath = ""
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
	if c.decompressProvider == nil {
		// TODO specify error
		return ErrNotOpenned
	}
	return c.decompressProvider.Walk(walkFn)
}

// Retrieve file relative to "provider" to destination
func (c *Gitlab) Retrieve(src string, dest string) error {
	return c.decompressProvider.Retrieve(src, dest)
}
