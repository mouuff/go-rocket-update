package main

import (
	"log"
	"runtime"
	"sync"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

// This example shows how you can run the update in background

func main() {

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries_" + runtime.GOOS + ".zip",
		},
		ExecutableName: "go-rocket-update-example",
		Version:        "v0.3.0",
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		updateStatus, err := u.Update()
		if err != nil {
			log.Println(err)
		}
		if updateStatus == updater.Updated {
			log.Println("Updated!")
		}

		wg.Done()
	}()
	log.Println(u.Version)
	wg.Wait()
}
