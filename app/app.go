package app

import (
	"atctest/atcoder"
	"errors"
	"flag"
	"fmt"
	"io"
)

const baseURL = "https://atcoder.jp"

const (
	exitCodeOK = iota
	exitCodeGetSamplesError
	exitCodeCheckError
)

type App struct {
	client  *atcoder.Client
	checker *atcoder.Checker

	outStream io.Writer
	errStream io.Writer
}

func New(args []string, outStream, errStream io.Writer) (*App, error) {
	flags := flag.NewFlagSet("atctest", flag.ContinueOnError)
	flags.SetOutput(errStream)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(errStream, helpMessage)
		flags.PrintDefaults()
	}

	var (
		contest string
		problem string
		command string
	)
	flags.StringVar(&contest, "contest", "", "contest you are challenging. e.g.) ABC051")
	flags.StringVar(&problem, "problem", "", "problem you are solving. e.g.) C")
	flags.StringVar(&command, "command", "", "command to execute your program. e.g.) 'python c.py'")
	if err := flags.Parse(args[1:]); err != nil {
		return nil, errors.New("ERROR: failed to parse flags")
	}

	if contest == "" {
		flags.Usage()
		return nil, errors.New("ERROR: specify the contest you are challenging. e.g.) ABC051")
	}
	if problem == "" {
		flags.Usage()
		return nil, errors.New("ERROR: specify the problem you are solving. e.g.) C")
	}
	if command == "" {
		flags.Usage()
		return nil, errors.New("ERROR: specify the command to execute your program. e.g.) 'python c.py'")
	}

	client, err := atcoder.NewClient(baseURL, contest, problem)
	if err != nil {
		return nil, err
	}

	checker := atcoder.NewChecker(command)

	return &App{
		client:  client,
		checker: checker,

		outStream: outStream,
		errStream: errStream,
	}, nil
}

func (a *App) Run() int {
	samples, err := a.client.GetSamples()
	if err != nil {
		_, _ = fmt.Fprintf(a.errStream, err.Error())
		return exitCodeGetSamplesError
	}

	if success := a.checker.Check(samples); !success {
		return exitCodeCheckError
	}

	return exitCodeOK
}

const helpMessage = `
atctest is a command line tool for AtCoder.
it checks if your program correctly solve the samples provided on the problem page.

EXAMPLE: 
$ atctest -contest ABC051 -problem C -command 'python c.py'

OPTION:
`
