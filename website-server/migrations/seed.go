package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/pkg/storer"
)

var envPath = flag.String("env", ".env.dev", "env path")

func main() {
	cfg := config.NewConfig(*envPath)
	s := storer.NewGcpStorer(cfg.App.GcpBucketName)
	err := s.DeleteFolder("testcases")
	if err != nil {
		fmt.Printf("s.DeleteFolder: %v", err)
	}

	entries, err := os.ReadDir("migrations/testcases")
	if err != nil {
		log.Fatalf("os.ReadDir: %v", err)
	}

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subEntries, err := os.ReadDir("migrations/testcases/" + entry.Name())
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if subEntry.IsDir() {
				continue
			}

			remoteDest := fmt.Sprintf("testcases/%s/%s", entry.Name(), subEntry.Name())
			localDest := "migrations/" + remoteDest

			f, err := os.Open(localDest)
			if err != nil {
				continue
			}

			wg.Add(1)
			sem <- struct{}{}
			go func(f *os.File, remoteDest string) {
				fmt.Println("start uploading")
				if _, err = s.Upload(f, remoteDest, true); err != nil {
					fmt.Printf("error: %v", err)
				}
				fmt.Println("\tstop uploading")
				wg.Done()
				<-sem
			}(f, remoteDest)
		}
	}
	wg.Wait()
}
