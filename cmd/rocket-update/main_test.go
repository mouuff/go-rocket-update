package main_test

import (
	"path/filepath"
	"testing"

	main "github.com/mouuff/go-rocket-update/cmd/rocket-update"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func TestKeyGen(t *testing.T) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	folderToSign := filepath.Join("testdata", "Allum1")
	privKeyPath := filepath.Join(tmpDir, "test_key")
	pubKeyPath := filepath.Join(tmpDir, "test_key.pub")

	err = main.RunSubCommands([]string{"keygen", "-name", privKeyPath})
	if err != nil {
		t.Fatal(err)
	}
	err = main.RunSubCommands([]string{"sign", "-path", folderToSign, "-key", privKeyPath})
	if err != nil {
		t.Fatal(err)
	}
	err = main.RunSubCommands([]string{"verify", "-path", folderToSign, "-pubkey", pubKeyPath})
	if err != nil {
		t.Fatal(err)
	}

}
