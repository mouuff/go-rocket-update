package provider_test

import (
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/pkg/provider"
)

func TestProviderSecure(t *testing.T) {

	pubStr := `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA4L9SG+I23Rr3rCv8MUWb
XEF2/SfOQBoYM80pKv4Q6CDICH5wmmWWfTchAEbkKHp3Hx9qP2/aR3mg1qLqWrXZ
hrYMFKvQCuu+3qGU+KxZfb9qtwIyMSdQdw3yLo7dw/1Kvwcowu6m6OTt/0buF7mr
UQwDM3EhyOCfKDBOJ8+WbAevvpnCI4fcGT9mrEZ6S2wgIx+6Io5Bbre7+iQ7l9RO
kE1iKITwXbCtpU3VDRkjtIQ3O0LdGMifJ+DNNIogxiAUGpYZVYVXummX9vBU1whr
bd2lf9PUB0HLEZPtMFgR3FtQpGqGSi2LXH7C7YSQDHdX/VWRb5kfmjJA/M57tD55
oe2F3Cxeqi3TNQ/8d+9Ta/8TUCrthyd0dPe0F98HlkyqFh+aIWuaK9mdhh/rzOk8
on97gotr1tYOdFTK09KHnYrSdCmqfvByOyQnHYvzqN6qwmEI3ufY/i/KTe6QxtWL
J73aa5rBKLCE/TYT53R2ZhFCXYTPhzwl5LkwYKdKIQ55Z9TfYhjbArjCMTz19Akl
6qjQcxAmGK2mY2odUkHD/Rfqs02fQoQdnRoi3qhkDW4fDYfcqTnkfTkcEvvJJTuI
ylMx5Dy8lG6/J7zKWkV6S7h3+K11dWZn7toQVyVU3M2GpEng3b74Pp3Ma7ymoM8J
SZ5Uz050oR/PoLaSx3xdjFMCAwEAAQ==
-----END RSA PUBLIC KEY-----`
	p := &provider.Secure{
		BackendProvider: &provider.Zip{Path: filepath.Join("testdata", "Allum1Signed-v1.0.0.zip")},
		PublicKeyPEM:    []byte(pubStr),
	}
	if err := p.Open(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	err := ProviderTestWalkAndRetrieve(p)
	if err != nil {
		t.Fatal(err)
	}

	// Changing public key for bad one
	pubStr = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA4L9SG+I23Rr3rCv8MUWb
XEF2/SfOQBoYM80pKv4Q6CDICH5wmmWWfTchAEbkKHp3Hx9qP2/aR3mg1qLqWrXZ
hrYMFKvQCuu+3qGU+KxZfb9qtwIyMSdQdw3yLo7dw/1Kvwcowu6m6OTt/0buF7mr
UQwDM3EhyOCfKDBOJ8+WbAevvpnCI4fcGT9mrEZ6S2wgIx+6Io5Bbre7+iQ7l9RO
kE1iKITwXbCtpU3VDRkjtIQ3O0LdGMifJ+DNNIogxiAUGpYZVYVXummX9vBU1whr
bd2lf9PUB0HLEZPtMFgR3FtQpGqGSi2LXH7C7YSQDHdX/VWRb5kfmjJA/M57tD55
oe2F3Cxeqi3TNQ/8d+9Ta/8TUCrthyd0dPe0F98HlkyqFh+aIWuaK9mdhh/rzOk8
on97gotr1tYOdFTK09KHnYrSdCmqfvByOyQnHYvzqN6qwmEI3ufY/i/KTe6QxtWL
J73aa5rBKLCE/TYT53R2ZhFCXYTPhzwl5LkwYKdKIQ55Z9TfYhjbArjCMTz19Akl
6qjQcxAmGK2mY2odUkHD/Rfqs02fQoQdnRoi3qhkDW4fDYfcqTnkfTkcEvvJJTuI
ylMx5Dy8lG6/J7zKWkV6S7h3+K11dWZn7toQVyVU3M2GpEng3b74Pp3Ma7zmoM8J
SZ5Uz050oR/PoLaSx3xdjFMCAwEAAQ==
-----END RSA PUBLIC KEY-----`
	p.PublicKey, err = crypto.ParsePemPublicKey([]byte(pubStr))
	if err != nil {
		t.Fatal(err)
	}
	err = ProviderTestWalkAndRetrieve(p)
	if err == nil {
		t.Fatal("ProviderTestWalkAndRetrieve shouldn't work with a bad public key")
	}

	badProvider := &provider.Secure{
		BackendProvider: &provider.Zip{Path: filepath.Join("testdata", "doesnotexist.zip")},
		PublicKeyPEM:    []byte(pubStr),
	}
	err = ProviderTestUnavailable(badProvider)
	if err != nil {
		t.Fatal(err)
	}
}
