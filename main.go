package main

import (
	"fmt"
	"log"

	"github.com/mouuff/easy-update/provider"
	"github.com/mouuff/easy-update/updater"
)

func main() {

	p := &provider.Github{
		RepositoryURL: "https://github.com/QuailTeam/cpp_indie_studio",
		ZipName:       "indie.zip",
	}
	fmt.Println(p.Open())

	fmt.Println(updater.GetPlatformName())

	u := &updater.Updater{
		Provider:   &provider.Local{Path: "testdata"},
		BinaryName: "main",
		Version:    "1.0",
	}
	log.Println(u.Version)
	err := u.Run()
	if err != nil {
		log.Fatal(err)
	}
}
