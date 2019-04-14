package main

import (
	"atctest/atcoder"
	"flag"
	"log"
)

func main() {
	var (
		contest string
		problem string
		command string
	)

	flag.StringVar(&contest, "contest", "", "contest. e.g) ABC051")
	flag.StringVar(&problem, "problem", "", "problem. e.g) A")
	flag.StringVar(&command, "command", "", "command. e.g) 'python a.py'")
	flag.Parse()

	samples, err := atcoder.GetSamples(contest, problem)
	if err != nil {
		log.Fatal(err)
	}

	atcoder.Check(samples, command)
}
