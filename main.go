package main

import (
	"fmt"
	"os"

	"github.com/mui87/atctest/app"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

func main() {
	a, err := app.New(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitCodeErr)
	}

	if err := a.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitCodeErr)
	}

	os.Exit(exitCodeOK)
}
