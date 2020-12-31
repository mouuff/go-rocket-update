package crypto_test

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/crypto"
)

func TestSignAndVerifyFile(t *testing.T) {
	fileA := filepath.Join("testdata", "small.txt")
	fileB := filepath.Join("testdata", "bin.txt")

	privA, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}
	privB, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}

	signatureA, err := crypto.GetFileSignature(privA, fileA)
	if err != nil {
		t.Error(err)
	}
	signatureB, err := crypto.GetFileSignature(privA, fileB)
	if err != nil {
		t.Error(err)
	}

	err = crypto.VerifyFileSignature(&privA.PublicKey, signatureA, fileA)
	if err != nil {
		t.Error(err)
	}
	err = crypto.VerifyFileSignature(&privB.PublicKey, signatureA, fileA)
	if err == nil {
		t.Error("fileA should not verify with privB")
	}
	err = crypto.VerifyFileSignature(&privB.PublicKey, signatureB, fileA)
	if err == nil {
		t.Error("fileA should not verify with signatureB")
	}
	err = crypto.VerifyFileSignature(&privA.PublicKey, signatureB, fileB)
	if err != nil {
		t.Error(err)
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
		t.Error(err)
	}
	err = verifyChecksumFileSHA256(fileB, "0596cc0127626799289943332342b56787cc589b1811f3b5a1fa108938765fa0")
	if err != nil {
		t.Error(err)
	}

}
