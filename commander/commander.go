package commander

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Commander interface {
	Run(rawCommand, stdin string) (string, error)
}

type External struct{}

func NewExternal() *External {
	return &External{}
}

func (e *External) Run(rawCommand, stdin string) (string, error) {
	var errBuf bytes.Buffer

	cmd := NewCommand(rawCommand)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stderr = &errBuf

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), errBuf.String())
	}
	return string(out), nil
}

func NewCommand(rawCommand string) *exec.Cmd {
	fields := strings.Fields(rawCommand)
	name := fields[0]
	args := fields[1:]

	return exec.Command(name, args...)
}
