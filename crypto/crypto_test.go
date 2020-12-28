package crypto_test

import (
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/crypto"
)

func TestSignAndVerifyFile(t *testing.T) {
	fileA := filepath.Join("testdata", "small.txt")
	fileB := filepath.Join("testdata", "bin.txt")

	privA, err := crypto.RandomPrivateKey()
	if err != nil {
		t.Error(err)
	}
	privB, err := crypto.RandomPrivateKey()
	if err != nil {
		t.Error(err)
	}

	signatureA, err := crypto.GetSignature(privA, fileA)
	if err != nil {
		t.Error(err)
	}
	signatureB, err := crypto.GetSignature(privA, fileB)
	if err != nil {
		t.Error(err)
	}

	err = crypto.VerifySignature(&privA.PublicKey, signatureA, fileA)
	if err != nil {
		t.Error(err)
	}
	err = crypto.VerifySignature(&privB.PublicKey, signatureA, fileA)
	if err == nil {
		t.Error("fileA should not verify with privB")
	}
	err = crypto.VerifySignature(&privB.PublicKey, signatureB, fileA)
	if err == nil {
		t.Error("fileA should not verify with signatureB")
	}
	err = crypto.VerifySignature(&privA.PublicKey, signatureB, fileB)
	if err != nil {
		t.Error(err)
	}

}
