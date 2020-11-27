package main

import (
	"log"

	"github.com/aerfio/sdxo/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
