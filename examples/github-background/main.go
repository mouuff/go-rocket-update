package main

import (
	"log"
	"runtime"
	"sync"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

// This example shows how you can can the update in background

func main() {

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries_" + runtime.GOOS + ".zip",
		},
		BinaryName: "go-rocket-update-example",
		Version:    "v0.3.0",
	}

	var wg sync.WaitGroup

	canUpdate, err := u.CanUpdate()
	if err != nil {
		log.Println(err)
	} else if canUpdate {
		log.Println("Update found! Updating in background...")
		wg.Add(1)
		go func() {
			if err := u.Update(); err != nil {
				log.Println(err)
			}

			wg.Done()
		}()
	} else {
		log.Println("No update found")
	}
	log.Println(u.Version)
	wg.Wait()
}
