package main

import (
	"log"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

func main() {
	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries.zip",
		},
		BinaryName: "go-rocket-update-example",
		Version:    "v0.0",
	}
	log.Println(u.Version)
	err := u.Run()
	if err != nil {
		log.Fatal(err)
	}
}
