package main

import (
	"fmt"
	"log"

	"github.com/mouuff/easy-update/provider"
	"github.com/mouuff/easy-update/updater"
)

func main() {

	p := provider.NewProviderGithub("https://github.com/QuailTeam/cpp_indie_studio", "indie.zip")
	fmt.Println(p.Open())

	fmt.Println(updater.GetPlatformName())

	u := &updater.Updater{
		Provider:   provider.NewProviderLocal("testdata"),
		BinaryName: "main",
		Version:    "1.0",
	}
	log.Println(u.Version)
	err := u.Run()
	if err != nil {
		log.Fatal(err)
	}
}
