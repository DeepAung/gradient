package main

import (
	"log"

	"github.com/DeepAung/gradient/cmd/website/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
