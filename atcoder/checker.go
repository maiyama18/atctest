package atcoder

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Checker struct {
	command string
}

func NewChecker(command string) *Checker {
	return &Checker{command: command}
}

func (c *Checker) Check(samples []Sample) (bool, error) {
	fields := strings.Fields(c.command)
	cm := fields[0]
	args := fields[1:]

	for i, sample := range samples {
		cmd := exec.Command(cm, args...)
		cmd.Stdin = strings.NewReader(sample.Input)

		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		result := "x"
		if string(out) == sample.Output {
			result = "o"
		}
		fmt.Printf("sample %d -> %s\n", i+1, result)
		fmt.Print("expected:", sample.Output)
		fmt.Print("actual:", string(out))
		fmt.Println()
	}

	return true, nil
}
