package crypto_test

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func TestSignAndVerifyFile(t *testing.T) {
	fileA := filepath.Join("testdata", "small.txt")
	fileB := filepath.Join("testdata", "bin.txt")

	privA, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	privB, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	signatureA, err := crypto.GetFileSignature(privA, fileA)
	if err != nil {
		t.Fatal(err)
	}
	signatureB, err := crypto.GetFileSignature(privA, fileB)
	if err != nil {
		t.Fatal(err)
	}

	err = crypto.VerifyFileSignature(&privA.PublicKey, signatureA, fileA)
	if err != nil {
		t.Fatal(err)
	}
	err = crypto.VerifyFileSignature(&privB.PublicKey, signatureA, fileA)
	if err == nil {
		t.Fatal("fileA should not verify with privB")
	}
	err = crypto.VerifyFileSignature(&privB.PublicKey, signatureB, fileA)
	if err == nil {
		t.Fatal("fileA should not verify with signatureB")
	}
	err = crypto.VerifyFileSignature(&privA.PublicKey, signatureB, fileB)
	if err != nil {
		t.Fatal(err)
	}

	// Error paths
	_, err = crypto.GetFileSignature(privA, filepath.Join("testdata", "doesnotexist.txt"))
	if err == nil {
		t.Fatal("Should return error when file does not exist")
	}
}

func TestVerifyFile(t *testing.T) {
	fileA := filepath.Join("testdata", "small.txt")

	pub, err := crypto.ParsePemPublicKey([]byte(`-----BEGIN RSA PUBLIC KEY-----
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
-----END RSA PUBLIC KEY-----`))
	if err != nil {
		t.Fatal(err)
	}
	signatureHex := `6aec38644a796635fa5d956e9178649a7160b94cb089bc90b7bac867cf4e586e97c0efcffdc52f2279da44132323ab1487e3f4801654b71001a8e5965f4e0c542ae2653bc50ec2a235ad86f11f1a879474f9be9df569d430406663fd7359f526d84283f94b07bd96bb3a79f3a24e94cc409631c25e19dbff7f9070320841d69ec6b7512d8d401b4a666466e5cb880f2e28ca09e997890b31b5c73c54c33c9a88a8367cb74612d9dcc3d59c14a64314c5eff18d849c6042b1dd02a49c15ae6a46f0130cef3d36305f34767657f9ea2ce81a6ed6d27318f32aa52d04d47994abcc97c586510b74ed1cfe23b41c5cf2eb08e8ebdbc7f6ed447142a37426673d48426101a9e422489a4d3f0d3325a0904c867c255d74211ec6beec66c73cb50953474ce03e571a128e2b6c049a08f09d8009afd926a1745c1231df3fe0ebe5374e1cd62738e0b822f0aab4ec2cf813bc6b63fa6de75e6adb1efc1001fbda2d3d8a1335109d2d82f33fc927f2ea93bbf52d2ac99834adc4b8e075dc78a74a01c90ef705e58436640fc19157e568ec20892c310fd3b3522b8349d69259e6bca7ecc4975ff01240a3591e2b7ed97a8df8754d36f9785a324bbad9525043a750f492ff42d8fa72f24749d3437ee5411aa4739696e9ab836fb888c3b9a1e6434f506a4b31431d185557525ee9323b9e8cd2034e046c5c75f1b5b252bbc0b49cc003f5cd69`
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		t.Fatal(err)
	}
	err = crypto.VerifyFileSignature(pub, signature, fileA)
	if err != nil {
		t.Fatal(err)
	}

	// Error paths
	err = crypto.VerifyFileSignature(pub, signature, filepath.Join("testdata", "doesnotexist.txt"))
	if err == nil {
		t.Fatal("Should return error when file does not exist.")
	}

	_, err = crypto.ParsePemPublicKey([]byte(`-----BEGIN PUBLIC KEY-----
	MCowBQYDK2VwAyEABhjHE6AOa33q2JGlVk9OjICRp2S6d9nUJh0Xr6PUego=
	-----END PUBLIC KEY-----`))
	if err == nil {
		t.Fatal("Should return error when type of the public key is not RSA.")
	}
}

func verifyChecksumFileSHA256(path string, expectedHexChecksum string) error {
	checksum, err := crypto.ChecksumFileSHA256(path)
	hexChecksum := hex.EncodeToString(checksum)
	if err != nil {
		return err
	}
	if hexChecksum != expectedHexChecksum {
		return fmt.Errorf("TestChecksumFileSHA256 file %s: %s != %s", path, hexChecksum, expectedHexChecksum)
	}
	return nil
}

func TestChecksumFileSHA256(t *testing.T) {
	fileA := filepath.Join("testdata", "small.txt")
	fileB := filepath.Join("testdata", "bin.txt")
	err := verifyChecksumFileSHA256(fileA, "f2a65cb3c3170bfe938f30e4dd592bfdd6c1b69b3a92046ef43b375d1eff669e")
	if err != nil {
		t.Fatal(err)
	}
	err = verifyChecksumFileSHA256(fileB, "0596cc0127626799289943332342b56787cc589b1811f3b5a1fa108938765fa0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportImport(t *testing.T) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	pub := &priv.PublicKey
	privExported := crypto.ExportPrivateKeyAsPem(priv)
	pubExported, err := crypto.ExportPublicKeyAsPem(pub)
	if err != nil {
		t.Fatal(err)
	}
	privImported, err := crypto.ParsePemPrivateKey(privExported)
	if err != nil {
		t.Fatal(err)
	}
	pubImported, err := crypto.ParsePemPublicKey(pubExported)
	if err != nil {
		t.Fatal(err)
	}
	if !pub.Equal(pubImported) {
		t.Error("pub != pubImported")
	}
	if !priv.Equal(privImported) {
		t.Error("priv != privImported")
	}

	_, err = crypto.ParsePemPublicKey(privExported)
	if err == nil {
		t.Error("ParsePemPublicKey should not work with public keys")
	}

	_, err = crypto.ParsePemPrivateKey(pubExported)
	if err == nil {
		t.Error("ParsePemPrivateKey should not work with private keys")
	}

	_, err = crypto.ParsePemPublicKey([]byte("Should not parse this"))
	if err == nil {
		t.Error("ParsePemPublicKey should not parse the input")
	}

	_, err = crypto.ParsePemPrivateKey([]byte("Should not parse this"))
	if err == nil {
		t.Error("ParsePemPrivateKey should not parse the input")
	}

	_, err = crypto.ExportPublicKeyAsPem(&rsa.PublicKey{})
	if err == nil {
		t.Error("ExportPublicKeyAsPem should not parse the input")
	}
}
