package main

import (
	"atctest/app"
	"log"
	"os"
)

func main() {
	a, err := app.New(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		log.Fatal(err)
		os.Exit(128)
	}

	os.Exit(a.Run())
}
