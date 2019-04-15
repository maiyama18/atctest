package commander

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Commander interface {
	Run(stdin string) (string, error)
}

type External struct {
	command *exec.Cmd
}

func NewExternal(rawCommand string) *External {
	fields := strings.Fields(rawCommand)
	name := fields[0]
	args := fields[1:]

	return &External{
		command: exec.Command(name, args...),
	}
}

func (e *External) Run(stdin string) (string, error) {
	var errBuf bytes.Buffer

	cmd := *e.command
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stderr = &errBuf

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), errBuf.String())
	}
	return string(out), nil
}
