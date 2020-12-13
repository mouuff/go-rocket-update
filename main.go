package main

import (
	"fmt"
	"log"

	"github.com/mouuff/easy-update/updater"
)

func main() {
	fmt.Println(updater.GetPlatformName())
	fmt.Println(updater.GetExecutable())
	fmt.Println(updater.GetPlatformName())
	fmt.Println(updater.GetExecutable())

	err := updater.ReplaceExecutableWith("mainold.exe")
	if err != nil {
		log.Fatal(err)
	}
}
