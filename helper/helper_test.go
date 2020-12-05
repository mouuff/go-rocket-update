package helper_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/mouuff/easy-update/helper"
)

func TestCopyFile(t *testing.T) {

	dir := os.TempDir()

	fmt.Print(dir)

	helper.CopyFile(path.Join("testdata", "small.exe"), path.Join(dir, "small2.txt"))
	t.Error(dir)
	//defer os.RemoveAll(dir)
}
