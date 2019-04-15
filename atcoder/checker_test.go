package atcoder

import (
	"os/exec"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name            string
		inputCommand    *exec.Cmd
		inputSamples    []Sample
		expectedSuccess bool
		expectedOutput  string
	}{
		{
			name:         "success",
			inputCommand: exec.Command("ruby", "c.rb"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
