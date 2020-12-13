package main

import (
	"fmt"

	"github.com/mouuff/easy-update/updater"
)

func main() {
	fmt.Println(updater.GetPlatformName())
	fmt.Println(updater.GetExecutable())
}
