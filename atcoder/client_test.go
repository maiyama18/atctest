package atcoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

const dummyBaseURL = "https://dummyatcoder.jp"

var dummyCacheDirPath = path.Join("testdata", "cache")

type mockInfo struct {
	path       string
	statusCode int
}

func TestClient_getProblemURL(t *testing.T) {
	tests := []struct {
		name string

		inputContest string
		inputProblem string

		mockInfos []mockInfo

		expectedProblemURL string
		expectedErrMsg     string
	}{
		{
			name:         "success-url_type1",
			inputContest: "abc120",
			inputProblem: "c",
			mockInfos: []mockInfo{
				{"/contests/abc120/tasks/abc120_c", http.StatusOK},
			},
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc120/tasks/abc120_c",
		},
		{
			name:         "success-url_type2",
			inputContest: "abc002",
			inputProblem: "c",
			mockInfos: []mockInfo{
				{"/contests/abc002/tasks/abc002_c", http.StatusNotFound},
				{"/contests/abc002/tasks/abc002_3", http.StatusOK},
			},
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc002/tasks/abc002_3",
		},
		{
			name:         "failure-nonexistent problem",
			inputContest: "xxx999",
			inputProblem: "x",
			mockInfos: []mockInfo{
				{"/contests/xxx999/tasks/xxx999_x", http.StatusNotFound},
				{"/contests/xxx999/tasks/xxx999_0", http.StatusNotFound},
			},
			expectedErrMsg: "could not find problem page",
		},
		{
			name:         "failure-nonexistent contest",
			inputContest: "xxx999",
			inputProblem: "c",
			mockInfos: []mockInfo{
				{"/contests/xxx999/tasks/xxx999_c", http.StatusNotFound},
				{"/contests/xxx999/tasks/xxx999_3", http.StatusNotFound},
			},
			expectedErrMsg: "could not find problem page",
		},
		{
			name:         "failure-nonexistent problem in existent contest",
			inputContest: "abc120",
			inputProblem: "x",
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

			c := &Client{baseURL: dummyBaseURL, contest: test.inputContest, problem: test.inputProblem}
			problemURL, err := c.getProblemURL()
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}
				if problemURL != test.expectedProblemURL {
					t.Fatalf("problem URL wrong. want='%s', got='%s'", test.expectedProblemURL, problemURL)
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

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string

		inputContest      string
		inputProblem      string
		inputUseCache     bool
		inputCacheDirPath string

		expectedContest string
		expectedProblem string
	}{
		{
			name:            "success",
			inputContest:    "ABC120",
			inputProblem:    "C",
			expectedContest: "abc120",
			expectedProblem: "c",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var outBuff, errBuff bytes.Buffer
			c := NewClient(dummyBaseURL, test.inputContest, test.inputProblem, test.inputUseCache, test.inputCacheDirPath, &outBuff, &errBuff)
			if c.contest != test.expectedContest {
				t.Fatalf("contest wrong. want='%s', got='%s'", test.expectedContest, c.contest)
			}
			if c.problem != test.expectedProblem {
				t.Fatalf("problem wrong. want='%s', got='%s'", test.expectedProblem, c.problem)
			}
		})
	}
}

func TestClient_constructSamples(t *testing.T) {
	tests := []struct {
		name string

		inputElements map[string]string

		expectedSamples []Sample
		expectedErrMsg  string
	}{
		{
			name: "success-multiple_samples",
			inputElements: map[string]string{
				"入力例 1": "1 3 5\n",
				"出力例 1": "9\n",
				"入力例 2": "2 4\n",
				"出力例 2": "6\n",
			},
			expectedSamples: []Sample{
				{Input: "1 3 5\n", Output: "9\n"},
				{Input: "2 4\n", Output: "6\n"},
			},
		},
		{
			name: "success-single_sample",
			inputElements: map[string]string{
				"入力例": "1 3 5\n",
				"出力例": "9\n",
			},
			expectedSamples: []Sample{
				{Input: "1 3 5\n", Output: "9\n"},
			},
		},
		{
			name:           "failure-no_element",
			inputElements:  map[string]string{},
			expectedErrMsg: "no sample",
		},
		{
			name: "failure-numbers_of_elements_odd",
			inputElements: map[string]string{
				"入力例 1": "1 3 5\n",
				"出力例 1": "9\n",
				"入力例 2": "2 4\n",
			},
			expectedErrMsg: "number of sample elements should be even",
		},
		{
			name: "failure-index_of_入力例_wrong",
			inputElements: map[string]string{
				"入力例 2": "1 3 5\n",
				"出力例 1": "9\n",
			},
			expectedErrMsg: "could not find '入力例 1'",
		},
		{
			name: "failure-index_of_出力例_wrong",
			inputElements: map[string]string{
				"入力例 1": "1 3 5\n",
				"出力例 2": "9\n",
			},
			expectedErrMsg: "could not find '出力例 1'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &Client{}
			samples, err := c.constructSamples(test.inputElements)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err.Error())
				}
				if len(samples) != len(test.expectedSamples) {
					t.Fatalf("length of samples wrong. want=%d, got=%d", len(test.expectedSamples), len(samples))
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
					t.Fatalf("error message %q is expected to contain %q", err.Error(), test.expectedErrMsg)
				}
			}
		})
	}
}

func TestClient_fetchSampleElements(t *testing.T) {
	tests := []struct {
		name string

		inputProblemURL string

		mockRequestPath string
		mockStatusCode  int
		mockHTMLFile    string

		expectedSampleElements map[string]string
		expectedErrMsg         string
	}{
		{
			name:            "success",
			inputProblemURL: dummyBaseURL + "/contests/abc124/tasks/abc124_b",
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/abc124/tasks/abc124_b",
			mockHTMLFile:    "abc124b.html",
			expectedSampleElements: map[string]string{
				"入力例 1": strings.Join([]string{
					"4",
					"6 5 6 8",
					"",
				}, "\n"),
				"出力例 1": "3\n",
				"入力例 2": strings.Join([]string{
					"5",
					"4 5 3 5 4",
					"",
				}, "\n"),
				"出力例 2": "3\n",
				"入力例 3": strings.Join([]string{
					"5",
					"9 5 6 8 4",
					"",
				}, "\n"),
				"出力例 3": "1\n",
			},
		},
		{
			name:            "success-old_DOM_structure",
			inputProblemURL: dummyBaseURL + "/contests/abc002/tasks/abc002_2",
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/abc002/tasks/abc002_2",
			mockHTMLFile:    "abc002c.html",
			expectedSampleElements: map[string]string{
				"入力例 1": strings.Join([]string{
					"1 0 3 0 2 5",
					"",
				}, "\n"),
				"出力例 1": "5.0\n",
				"入力例 2": strings.Join([]string{
					"-1 -2 3 4 5 6",
					"",
				}, "\n"),
				"出力例 2": "2.0\n",
				"入力例 3": strings.Join([]string{
					"298 520 903 520 4 663",
					"",
				}, "\n"),
				"出力例 3": "43257.5\n",
			},
		},
		{
			name:            "success-only one sample",
			inputProblemURL: dummyBaseURL + "/contests/kupc2015/tasks/kupc2015_a",
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/kupc2015/tasks/kupc2015_a",
			mockHTMLFile:    "kupc2015a.html",
			expectedSampleElements: map[string]string{
				"入力例": strings.Join([]string{
					"3",
					"higashikyoto",
					"kupconsitetokyotokyoto",
					"goodluckandhavefun",
					"",
				}, "\n"),
				"出力例": strings.Join([]string{
					"1",
					"2",
					"0",
					"",
				}, "\n"),
			},
		},
		{
			name:            "failure-nonexistent_problem_URL",
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

			c := &Client{}
			sampleElements, err := c.fetchSampleElements(test.inputProblemURL)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}
				if len(sampleElements) != len(test.expectedSampleElements) {
					t.Fatalf("size of samples wrong. want=%d, got=%d", len(test.expectedSampleElements), len(sampleElements))
				}
				for key, expected := range test.expectedSampleElements {
					actual, ok := sampleElements[key]
					if !ok {
						t.Fatalf("sample element for key %q does not exist", key)
					}
					if actual != expected {
						t.Fatalf("sample element for key %q wrong. want=%+v, got=%+v", key, expected, actual)
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

func TestClient_cacheSamples(t *testing.T) {
	tests := []struct {
		name string

		inputContest      string
		inputProblem      string
		inputCacheDirPath string
		inputSamples      []Sample

		expectedErrMsg string
	}{
		{
			name:              "success",
			inputContest:      "abc120",
			inputProblem:      "c",
			inputCacheDirPath: dummyCacheDirPath,
			inputSamples: []Sample{
				{Input: "input1\n", Output: "output1\n"},
				{Input: "input2\n", Output: "output2\n"},
			},
		},
		{
			name:              "success-cache_dir_not_exist",
			inputContest:      "abc120",
			inputProblem:      "c",
			inputCacheDirPath: path.Join("testdata", "new_cache_dir"),
			inputSamples: []Sample{
				{Input: "input1\n", Output: "output1\n"},
				{Input: "input2\n", Output: "output2\n"},
			},
		},
		{
			name:              "failed-cache_dir_failed_no_permission",
			inputContest:      "abc120",
			inputProblem:      "c",
			inputCacheDirPath: path.Join("/sys", "new_cache_dir"),
			inputSamples: []Sample{
				{Input: "input1\n", Output: "output1\n"},
				{Input: "input2\n", Output: "output2\n"},
			},
			expectedErrMsg: "permission denied",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := os.MkdirAll(dummyCacheDirPath, 0777); err != nil {
				t.Fatal("failed to create dummy cache dir")
			}
			defer func() {
				if err := os.RemoveAll(dummyCacheDirPath); err != nil {
					t.Fatalf("failed to remove dummy cache dir after test: %s", dummyCacheDirPath)
				}
				if test.inputCacheDirPath != dummyCacheDirPath {
					if err := os.RemoveAll(test.inputCacheDirPath); err != nil {
						t.Fatalf("failed to remove cache dir after test: %s", test.inputCacheDirPath)
					}
				}
			}()

			c := &Client{contest: test.inputContest, problem: test.inputProblem, cacheDirPath: test.inputCacheDirPath}
			err := c.cacheSamples(test.inputSamples)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}

				filename := fmt.Sprintf("%s-%s.json", c.contest, c.problem)
				b, err := ioutil.ReadFile(path.Join(c.cacheDirPath, filename))
				if err != nil {
					t.Fatalf("failed to read cached file")
				}
				var samples []Sample
				if err := json.Unmarshal(b, &samples); err != nil {
					t.Fatalf("failed to unmarshal samples")
				}

				if len(samples) != len(test.inputSamples) {
					t.Fatalf("length of samples wrong. want=%d, got=%d", len(test.inputSamples), len(samples))
				}
				for i, expected := range test.inputSamples {
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
