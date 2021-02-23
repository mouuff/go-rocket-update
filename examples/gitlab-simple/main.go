package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

func main() {

	// Example project here: https://gitlab.com/arnaudalies.py/go-rocket-update-example
	u := &updater.Updater{
		Provider: &provider.Gitlab{
			ProjectID:   24021648,
			ArchiveName: fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
		},
		ExecutableName: fmt.Sprintf("go-rocket-update-example_%s_%s", runtime.GOOS, runtime.GOARCH),
		Version:        "v0.0.0",
	}

	log.Println(u.Version)
	if _, err := u.Update(); err != nil {
		log.Println(err)
	}
}
