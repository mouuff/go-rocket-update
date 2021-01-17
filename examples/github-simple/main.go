package main

import (
	"log"
	"runtime"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

func main() {

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries_" + runtime.GOOS + ".zip",
		},
		BinaryName: "go-rocket-update-example",
		Version:    "v0.3.0",
	}

	log.Println(u.Version)
	if err := u.Update(); err != nil {
		log.Println(err)
	}
}
