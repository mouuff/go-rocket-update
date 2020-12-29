package crypto_test

import (
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
