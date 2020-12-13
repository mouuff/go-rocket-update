package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/osext"
)

func GetPlatformName() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

func GetExecutable() (string, error) {
	execPath, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return execPath, nil
}

// ReplaceExecutableWith replaces the current executable with the one located at src
func ReplaceExecutableWith(src string) error {
	executable, err := GetExecutable()
	if err != nil {
		return err
	}
	tmpDir, err := ioutil.TempDir("", "updater")
	if err != nil {
		return err
	}
	// Here we move the current executable to a tmp dir, we do that because
	// on windows we must move the running executable to rewrite it
	renamedExecutable := filepath.Join(tmpDir, filepath.Base(executable))
	fmt.Println(renamedExecutable)
	fmt.Println(renamedExecutable)
	err = os.Rename(executable, renamedExecutable)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		os.Rename(renamedExecutable, executable)
		return err
	}

	err = ioutil.WriteFile(executable, content, 0755)
	if err != nil {
		os.Rename(renamedExecutable, executable)
		return err
	}
	return nil
}
