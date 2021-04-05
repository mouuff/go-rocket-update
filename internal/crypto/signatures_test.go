package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func TestSignatures(t *testing.T) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	root := "testdata"
	signatures, err := crypto.GetFolderSignatures(priv, root)
	if err != nil {
		t.Fatal(err)
	}
	unverifiedFiles, err := signatures.VerifyFolder(&priv.PublicKey, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(unverifiedFiles) > 0 {
		t.Fatal("All files should be verified")
	}

	signaturesPath := filepath.Join(tmpDir, "signaturestest.json")
	err = crypto.WriteSignaturesToJSON(signaturesPath, signatures)
	if err != nil {
		t.Fatal(err)
	}

	signatures.Remove("bin.txt")
	unverifiedFiles, err = signatures.VerifyFolder(&priv.PublicKey, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(unverifiedFiles) < 1 {
		t.Fatal("bin.txt should not be verified")
	}

	// Load previously saved signatures
	signatures, err = crypto.LoadSignaturesFromJSON(signaturesPath)
	if err != nil {
		t.Fatal(err)
	}
	// Second VerifyFolder with loaded signatures
	unverifiedFiles, err = signatures.VerifyFolder(&priv.PublicKey, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(unverifiedFiles) > 0 {
		t.Fatal("All files should be verified")
	}
	// We replace the signature of bin.txt with the one of small.txxt
	signatures.Remove("bin.txt")
	signatureSmallTXT, err := signatures.Get("small.txt")
	if err != nil {
		t.Fatal(err)
	}
	signatures.Add("bin.txt", signatureSmallTXT)
	unverifiedFiles, err = signatures.VerifyFolder(&priv.PublicKey, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(unverifiedFiles) < 1 {
		t.Fatal("bin.txt should not be verified")
	}
}
