package main

import (
	"fmt"
	"log"

	"github.com/mouuff/easy-update/provider"
	"github.com/mouuff/easy-update/updater"
)

func main() {

	fmt.Println(updater.GetPlatformName())

	version := "0.0"

	u := updater.NewUpdater(
		provider.NewProviderLocal("testdata"),
		"main",
		version,
	)
	log.Println(version)
	err := u.Run()
	if err != nil {
		log.Fatal(err)
	}
}
