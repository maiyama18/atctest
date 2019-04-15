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
		_, _ = fmt.Fprintf(c.outStream, "sample %d: ", i+1)
		if err != nil {
			successAll = false

			_, _ = color.New(color.FgRed).Fprintln(c.outStream, "ERROR")
			_, _ = fmt.Fprintln(c.outStream, err.Error())
		} else if success {
			_, _ = color.New(color.FgGreen).Fprintln(c.outStream, "SUCCESS")
		} else {
			successAll = false

			_, _ = color.New(color.FgRed).Fprintln(c.outStream, "FAILURE")
			_, _ = fmt.Fprintln(c.outStream, "input:")
			_, _ = fmt.Fprint(c.outStream, sample.Input)
			_, _ = fmt.Fprintln(c.outStream, "expected output:")
			_, _ = fmt.Fprint(c.outStream, sample.Output)
			_, _ = fmt.Fprintln(c.outStream, "actual output:")
			_, _ = fmt.Fprint(c.outStream, actual)
		}
		_, _ = fmt.Fprintln(c.outStream, "")
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
