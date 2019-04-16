package atcoder

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestChecker_Check(t *testing.T) {
	tests := []struct {
		name            string
		inputSamples    []Sample
		mockResults     []commandResult
		expectedSuccess bool
		expectedOutput  string
	}{
		{
			name: "success",
			inputSamples: []Sample{
				{Input: "0 1\n", Output: "1\n"},
				{Input: "1 2\n", Output: "3\n"},
			},
			mockResults: []commandResult{
				{output: "1\n", err: nil},
				{output: "3\n", err: nil},
			},
			expectedSuccess: true,
			expectedOutput:  "SUCCESS",
		},
		{
			name: "failure-all failed",
			inputSamples: []Sample{
				{Input: "0 1\n", Output: "1\n"},
				{Input: "1 2\n", Output: "3\n"},
			},
			mockResults: []commandResult{
				{output: "99\n", err: nil},
				{output: "99\n", err: nil},
			},
			expectedSuccess: false,
			expectedOutput:  "FAILURE\ninput:\n0 1",
		},
		{
			name: "failure-some failed",
			inputSamples: []Sample{
				{Input: "0 1\n", Output: "1\n"},
				{Input: "1 2\n", Output: "3\n"},
			},
			mockResults: []commandResult{
				{output: "1\n", err: nil},
				{output: "99\n", err: nil},
			},
			expectedSuccess: false,
			expectedOutput:  "FAILURE\ninput:\n1 2",
		},
		{
			name: "failure-some error",
			inputSamples: []Sample{
				{Input: "0 1\n", Output: "1\n"},
				{Input: "1 2\n", Output: "3\n"},
			},
			mockResults: []commandResult{
				{output: "1\n", err: nil},
				{output: "", err: errors.New("some error")},
			},
			expectedSuccess: false,
			expectedOutput:  "ERROR\nsome error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var outStream bytes.Buffer
			c := &Checker{
				command:   &testCommander{index: 0, results: test.mockResults},
				outStream: &outStream,
			}

			actualSuccess := c.Check(test.inputSamples)
			if actualSuccess != test.expectedSuccess {
				t.Fatalf("success wrong. want=%t, got=%t", test.expectedSuccess, actualSuccess)
			}
			if !strings.Contains(outStream.String(), test.expectedOutput) {
				t.Fatalf("expect '%s' to contain '%s'", outStream.String(), test.expectedOutput)
			}
		})
	}
}

type commandResult struct {
	output string
	err    error
}

type testCommander struct {
	index   int
	results []commandResult
}

func (t *testCommander) Run(stdin string) (string, error) {
	if t.index >= len(t.results) {
		panic("index of testCommander out of range")
	}
	result := t.results[t.index]
	t.index++

	return result.output, result.err
}
