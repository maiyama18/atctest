package atcoder

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

const dummyBaseURL = "https://dummyatcoder.jp"

type mockInfo struct {
	path       string
	statusCode int
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string

		inputContest      string
		inputProblem      string
		inputCacheDirPath string

		mockInfos []mockInfo

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

			c, err := NewClient(dummyBaseURL, test.inputContest, test.inputProblem, test.inputCacheDirPath)
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

// TODO: cache関連のテスト追加
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
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/abc124/tasks/abc124_b",
			mockHTMLFile:    "abc124b.html",
			expectedSamples: []Sample{
				{
					Input: strings.Join([]string{
						"4",
						"6 5 6 8",
						"",
					}, "\n"),
					Output: "3\n",
				},
				{
					Input: strings.Join([]string{
						"5",
						"4 5 3 5 4",
						"",
					}, "\n"),
					Output: "3\n",
				},
				{
					Input: strings.Join([]string{
						"5",
						"9 5 6 8 4",
						"",
					}, "\n"),
					Output: "1\n",
				},
			},
		},
		{
			name:            "success-old DOM structure",
			inputProblemURL: dummyBaseURL + "/contests/abc002/tasks/abc002_2",
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/abc002/tasks/abc002_2",
			mockHTMLFile:    "abc002c.html",
			expectedSamples: []Sample{
				{
					Input: strings.Join([]string{
						"1 0 3 0 2 5",
						"",
					}, "\n"),
					Output: "5.0\n",
				},
				{
					Input: strings.Join([]string{
						"-1 -2 3 4 5 6",
						"",
					}, "\n"),
					Output: "2.0\n",
				},
				{
					Input: strings.Join([]string{
						"298 520 903 520 4 663",
						"",
					}, "\n"),
					Output: "43257.5\n",
				},
			},
		},
		{
			name:            "success-only one sample",
			inputProblemURL: dummyBaseURL + "/contests/kupc2015/tasks/kupc2015_a",
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/kupc2015/tasks/kupc2015_a",
			mockHTMLFile:    "kupc2015a.html",
			expectedSamples: []Sample{
				{
					Input: strings.Join([]string{
						"3",
						"higashikyoto",
						"kupconsitetokyotokyoto",
						"goodluckandhavefun",
						"",
					}, "\n"),
					Output: strings.Join([]string{
						"1",
						"2",
						"0",
						"",
					}, "\n"),
				},
			},
		},
		{
			name:            "failure-nonexistent problem URL",
			inputProblemURL: dummyBaseURL + "/contests/xxx999/tasks/xxx999_x",
			mockStatusCode:  http.StatusNotFound,
			mockHTMLFile:    "xxx999x.html",
			mockRequestPath: "contests/xxx999/tasks/xxx999_x",
			expectedErrMsg:  "could not get HTML",
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
				AddHeader("Content-Type", "text/html").
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
					actual := samples[i]
					if actual != expected {
						t.Fatalf("%d-th sample wrong. want=%+v, got=%+v", i, expected, actual)
					}
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
