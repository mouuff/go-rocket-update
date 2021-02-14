package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

func main() {

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
		},
		ExecutableName: fmt.Sprintf("go-rocket-update-example_%s_%s", runtime.GOOS, runtime.GOARCH),
		Version:        "v0.0.0",
	}

	log.Println(u.Version)
	if _, err := u.Update(); err != nil {
		log.Println(err)
	}
}
