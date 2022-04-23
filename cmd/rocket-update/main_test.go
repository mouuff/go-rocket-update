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
	err := main.RunSubCommand([]string{"verify", "-path", folder, "-publicKey", publicKeyPath})
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
	defer os.RemoveAll(tmpDir)

	folder := filepath.Join(tmpDir, "fakepackage")
	if err := CreateFakePackage(folder); err != nil {
		t.Fatal(err)
	}
	privateKeyPath := filepath.Join(tmpDir, "test_key")
	publicKeyPath := filepath.Join(tmpDir, "test_key.pub")

	if err = keyGen(privateKeyPath); err != nil {
		t.Fatal(err)
	}

	if err = signFolder(folder, privateKeyPath+"doesnotexist"); err == nil {
		t.Fatal("signFolder shouldn't work if private key does not exist")
	}
	if err = signFolder(folder, filepath.Join("testdata", "file.jpeg")); err == nil {
		t.Fatal("signFolder shouldn't work if file is not a private key")
	}
	if err = signFolder(folder+"x", privateKeyPath); err == nil {
		t.Fatal("signFolder shouldn't work if folder does not exists")
	}
	if err = signFolder(folder, privateKeyPath); err != nil {
		t.Fatal(err)
	}
	if err = verifyFolder(folder, publicKeyPath); err != nil {
		t.Fatal(err)
	}
	if err = verifyFolder(folder+"x", publicKeyPath); err == nil {
		t.Fatal("verifyFolder shouldn't not work if the path is wrong")
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

	if err = verifyFolder(folder, publicKeyPath+"doesnotexist"); err == nil {
		t.Fatal("Folder shouldn't be verified if public key doesn't exist")
	}

	if err = verifyFolder(folder, filepath.Join("testdata", "file.jpeg")); err == nil {
		t.Fatal("Folder shouldn't be verified when file isn't a public key")
	}

	err = main.RunSubCommand([]string{})
	if err == nil {
		t.Fatal("Should return an error")
	}
	err = main.RunSubCommand([]string{"Test"})
	if err == nil {
		t.Fatal("Command shouldn't exist")
	}

	if err = keyGen(privateKeyPath); err == nil {
		t.Fatal("Both keys already exists")
	}
	os.Remove(privateKeyPath)
	if err = keyGen(privateKeyPath); err == nil {
		t.Fatal("public key already exists")
	}
	os.Remove(publicKeyPath)
	if err = keyGen(privateKeyPath); err != nil {
		t.Fatal(err)
	}
	os.Remove(publicKeyPath)
	if err = keyGen(privateKeyPath); err == nil {
		t.Fatal("private key already exists")
	}
}
