package crypto_test

import (
	"testing"

	"github.com/mouuff/go-rocket-update/internal/crypto"
)

func TestSignatures(t *testing.T) {
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}
	root := "testdata"
	signatures, err := crypto.GetFolderSignatures(priv, root)
	if err != nil {
		t.Error(err)
	}
	unverifiedFiles, err := signatures.VerifyFolder(&priv.PublicKey, root)
	if err != nil {
		t.Error(err)
	}
	if len(unverifiedFiles) > 0 {
		t.Error("All files should be verified")
	}
	signatures.Remove("bin.txt")
	unverifiedFiles, err = signatures.VerifyFolder(&priv.PublicKey, root)
	if err != nil {
		t.Error(err)
	}
	if len(unverifiedFiles) < 1 {
		t.Error("bin.txt should not be verified")
	}
}
