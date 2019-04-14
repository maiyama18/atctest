package atcoder

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

type Checker struct {
	command   *exec.Cmd
	outStream io.Writer
	errStream io.Writer
}

func NewChecker(rawCommand string, outStream, errStream io.Writer) *Checker {
	fields := strings.Fields(rawCommand)
	cm := fields[0]
	args := fields[1:]

	return &Checker{
		command:   exec.Command(cm, args...),
		outStream: outStream,
		errStream: errStream,
	}
}

// TODO: use checker's out/errStreams
func (c *Checker) Check(samples []Sample) bool {
	successAll := true
	for i, sample := range samples {
		success, actual, err := c.checkOne(sample)
		fmt.Printf("sample %d: ", i+1)
		if err != nil {
			successAll = false

			color.Red("ERROR\n")
			fmt.Println(err)
		} else if success {
			color.Green("SUCCESS\n")
		} else {
			successAll = false

			color.Red("FAILURE\n")
			fmt.Println("input:")
			fmt.Print(sample.Input)
			fmt.Println("expected output:")
			fmt.Print(sample.Output)
			fmt.Println("actual output:")
			fmt.Print(actual)
		}
		fmt.Println()
	}

	return successAll
}

func (c *Checker) checkOne(sample Sample) (bool, string, error) {
	var errBuf bytes.Buffer

	cmd := *c.command
	cmd.Stdin = strings.NewReader(sample.Input)
	cmd.Stderr = &errBuf

	out, err := cmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("%s: %s", err.Error(), errBuf.String())
	}

	return string(out) == sample.Output, string(out), nil
}
