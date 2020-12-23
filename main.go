package main

import (
	"fmt"
	"log"

	"github.com/mouuff/go-rocket-update/provider"
	"github.com/mouuff/go-rocket-update/updater"
)

func main() {

	fmt.Println(updater.GetPlatformName())

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries.zip",
		},
		BinaryName: "rocket",
		Version:    "1.1",
	}
	log.Println(u.Version)
	err := u.Run()
	if err != nil {
		log.Fatal(err)
	}
}
