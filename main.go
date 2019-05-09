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
	os.Exit(run())
}

func run() int {
	a, err := app.New(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[ERROR] "+err.Error())
		return exitCodeErr
	}

	if err := a.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[ERROR] "+err.Error())
		return exitCodeErr
	}

	return exitCodeOK
}
