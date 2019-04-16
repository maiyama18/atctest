package atcoder

import (
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"
)

const dummyBaseURL = "https://dummyatcoder.jp"

type mockInfo struct {
	path       string
	statusCode int
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name               string
		inputContest       string
		inputProblem       string
		mockInfos          []mockInfo
		expectedProblemURL string
		expectedErrMsg     string
	}{
		{
			name:         "success-url_type1",
			inputContest: "ABC120",
			inputProblem: "C",
			mockInfos: []mockInfo{
				{"/contests/abc120/tasks/abc120_c", http.StatusOK},
			},
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc120/tasks/abc120_c",
		},
		{
			name:         "success-url_type2",
			inputContest: "ABC002",
			inputProblem: "C",
			mockInfos: []mockInfo{
				{"/contests/abc002/tasks/abc002_c", http.StatusNotFound},
				{"/contests/abc002/tasks/abc002_3", http.StatusOK},
			},
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc002/tasks/abc002_3",
		},
		{
			name:         "failure-nonexistent problem",
			inputContest: "XXX999",
			inputProblem: "X",
			mockInfos: []mockInfo{
				{"/contests/xxx999/tasks/xxx999_x", http.StatusNotFound},
				{"/contests/xxx999/tasks/xxx999_0", http.StatusNotFound},
			},
			expectedErrMsg: "could not find problem page",
		},
		{
			name:         "failure-nonexistent contest",
			inputContest: "XXX999",
			inputProblem: "C",
			mockInfos: []mockInfo{
				{"/contests/xxx999/tasks/xxx999_c", http.StatusNotFound},
				{"/contests/xxx999/tasks/xxx999_3", http.StatusNotFound},
			},
			expectedErrMsg: "could not find problem page",
		},
		{
			name:         "failure-nonexistent problem in existent contest",
			inputContest: "ABC120",
			inputProblem: "X",
			mockInfos: []mockInfo{
				{"/contests/abc120/tasks/abc120_x", http.StatusNotFound},
				{"/contests/abc120/tasks/abc120_0", http.StatusNotFound},
			},
			expectedErrMsg: "could not find problem page",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer gock.Off()
			for _, mockInfo := range test.mockInfos {
				gock.New(dummyBaseURL).
					Get(mockInfo.path).
					Reply(mockInfo.statusCode)
			}

			c, err := NewClient(dummyBaseURL, test.inputContest, test.inputProblem)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}
				if c.problemURL != test.expectedProblemURL {
					t.Fatalf("problem URL wrong. want='%s', got='%s'", test.expectedProblemURL, c.problemURL)
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

func TestClient_GetSamples(t *testing.T) {
	tests := []struct {
		name            string
		inputProblemURL string
		mockRequestPath string
		mockStatusCode  int
		mockHTMLFile    string
		expectedSamples []Sample
		expectedErrMsg  string
	}{
		{
			name:            "success",
			inputProblemURL: dummyBaseURL + "/contests/abc124/tasks/abc124_b",
			mockRequestPath: "/contests/abc124/tasks/abc124_b",
			mockHTMLFile:    "abc124b.html",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			html, err := ioutil.ReadFile(path.Join("testdata", test.mockHTMLFile))
			if err != nil {
				t.Fatal(err)
			}

			defer gock.Off()
			gock.New(dummyBaseURL).
				Get(test.mockRequestPath).
				Reply(test.mockStatusCode).
				BodyString(string(html))

			c := &Client{problemURL: test.inputProblemURL}
			samples, err := c.GetSamples()
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}
				if len(samples) != len(test.expectedSamples) {
					t.Fatalf("size of samples wrong. want=%d, got=%d", len(test.expectedSamples), len(samples))
				}
				for i, expected := range test.expectedSamples {

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
