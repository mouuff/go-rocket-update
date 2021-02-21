package main

import (
	"fmt"
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
			ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
		},
		ExecutableName: fmt.Sprintf("go-rocket-update-example_%s_%s", runtime.GOOS, runtime.GOARCH),
		Version:        "v0.0.0",
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
