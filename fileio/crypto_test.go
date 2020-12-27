package fileio_test

import (
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/fileio"
)

func TestSignAndVerifyFile(t *testing.T) {
	fileA := filepath.Join("testdata", "TempleOS.ISO")
	fileB := filepath.Join("testdata", "small.exe")

	privA, err := fileio.RandomPrivateKey()
	if err != nil {
		t.Error(err)
	}
	privB, err := fileio.RandomPrivateKey()
	if err != nil {
		t.Error(err)
	}

	signatureA, err := fileio.GetSignature(privA, fileA)
	if err != nil {
		t.Error(err)
	}
	signatureB, err := fileio.GetSignature(privA, fileB)
	if err != nil {
		t.Error(err)
	}

	err = fileio.VerifySignature(&privA.PublicKey, signatureA, fileA)
	if err != nil {
		t.Error(err)
	}
	err = fileio.VerifySignature(&privB.PublicKey, signatureA, fileA)
	if err == nil {
		t.Error("fileA should not verify with privB")
	}
	err = fileio.VerifySignature(&privB.PublicKey, signatureB, fileA)
	if err == nil {
		t.Error("fileA should not verify with signatureB")
	}
	err = fileio.VerifySignature(&privA.PublicKey, signatureB, fileB)
	if err != nil {
		t.Error(err)
	}

}
