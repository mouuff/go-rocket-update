package main_test

import (
	"os"
	"path/filepath"
	"testing"

	main "github.com/mouuff/go-rocket-update/cmd/rocket-update"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

func signFolder(folder string, privateKeyPath string) error {
	err := main.RunSubCommand([]string{"sign", "-path", folder, "-key", privateKeyPath})
	if err != nil {
		return err
	}
	return nil
}

func verifyFolder(folder string, publicKeyPath string) error {
	err := main.RunSubCommand([]string{"verify", "-path", folder, "-pubkey", publicKeyPath})
	if err != nil {
		return err
	}
	return nil
}

func keyGen(name string) error {
	err := main.RunSubCommand([]string{"keygen", "-name", name})
	if err != nil {
		return err
	}
	return nil
}

func CreateFakePackage(folder string) error {
	subFolderOne := filepath.Join(folder, "subfolder1")
	os.Mkdir(folder, os.ModePerm)
	os.Mkdir(subFolderOne, os.ModePerm)

	filenameList := []string{"binary", "file.jpeg", "file.txt"}

	for _, filename := range filenameList {
		err := fileio.CopyFile(
			filepath.Join("testdata", filename),
			filepath.Join(folder, filename))
		if err != nil {
			return err
		}
		err = fileio.CopyFile(
			filepath.Join("testdata", filename),
			filepath.Join(subFolderOne, filename))
		if err != nil {
			return err
		}
	}

	return nil
}

func TestMain(t *testing.T) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	folder := filepath.Join(tmpDir, "fakepackage")
	if err := CreateFakePackage(folder); err != nil {
		t.Fatal(err)
	}
	privateKeyPath := filepath.Join(tmpDir, "test_key")
	publicKeyPath := filepath.Join(tmpDir, "test_key.pub")

	if err = keyGen(privateKeyPath); err != nil {
		t.Fatal(err)
	}
	if err = signFolder(folder, privateKeyPath); err != nil {
		t.Fatal(err)
	}
	if err = verifyFolder(folder, publicKeyPath); err != nil {
		t.Fatal(err)
	}

	// Adding a file (which is not going to be verified)
	err = fileio.CopyFile(
		filepath.Join("testdata", "file.jpeg"),
		filepath.Join(folder, "test2.jpeg"))
	if err != nil {
		t.Fatal(err)
	}

	if err = verifyFolder(folder, publicKeyPath); err == nil {
		t.Fatal("Folder shouldn't be verified")
	}

}
