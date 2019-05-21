package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name               string
		inputArgs          []string
		expectedContestURL string
		expectedErrMsg     string
	}{
		{
			name:               "success",
			inputArgs:          strings.Fields("atctest -contest ABC051 -problem C -command 'python c.py'"),
			expectedContestURL: "https://atcoder.jp/contests/abc051",
		},
		{
			name:               "success-with url_old",
			inputArgs:          strings.Fields("atctest -url 'https://abc051.contest.atcoder.jp/tasks/abc051_c'"),
			expectedContestURL: "https://abc051.contest.atcoder.jp",
		},
		{
			name:               "success-with url_new",
			inputArgs:          strings.Fields("atctest -url 'https://atcoder.jp/contests/abc051/tasks/abc051_c'"),
			expectedContestURL: "https://atcoder.jp/contests/abc051",
		},
		{
			name:           "failure-unknown option exists",
			inputArgs:      strings.Fields("atctest -hello world -problem C -command 'python c.py'"),
			expectedErrMsg: "failed to parse flags",
		},
		{
			name:           "failure-contest option missing",
			inputArgs:      strings.Fields("atctest -problem C -command 'python c.py'"),
			expectedErrMsg: "specify the contest",
		},
		{
			name:           "failure-problem option missing",
			inputArgs:      strings.Fields("atctest -contest ABC051 -command 'python c.py'"),
			expectedErrMsg: "specify the problem",
		},
		{
			name:           "failure-command option missing",
			inputArgs:      strings.Fields("atctest -contest ABC051 -problem C"),
			expectedErrMsg: "specify the command",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var outStream, errStream bytes.Buffer
			a, err := New(test.inputArgs, &outStream, &errStream)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}
				if a.contestURL != test.expectedContestURL {
					t.Fatalf("contestURL wrong. want=%s, got=%s", test.expectedContestURL, a.contestURL)
				}
			} else {
				if err == nil {
					t.Fatal("err should not be nil. got: nil")
				}
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Fatalf("expect '%s' to contain '%s'", err.Error(), test.expectedErrMsg)
				}
			}
		})
	}
}
