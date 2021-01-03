package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"

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
		Version:    "v0.1",
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err := u.Update()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	log.Println(u.Version)
	fmt.Println("Hello world during update!")
	wg.Wait()
}
