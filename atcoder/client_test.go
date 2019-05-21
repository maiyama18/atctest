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

	"github.com/gocolly/colly"

	"gopkg.in/h2non/gock.v1"
)

const dummyBaseURL = "https://dummyatcoder.jp"

var dummyCacheDirPath = path.Join("testdata", "cache")

func TestClient_IsContestBeingHeld(t *testing.T) {
	tests := []struct {
		name string

		inputContestURL string

		mockRequestPath string
		mockStatusCode  int
		mockHTMLFile    string

		expected       bool
		expectedErrMsg string
	}{
		{
			name:            "success-apg4b_being_held",
			inputContestURL: dummyBaseURL + "/contests/APG4b",
			mockRequestPath: "/contests/APG4b",
			mockStatusCode:  http.StatusOK,
			mockHTMLFile:    "apg4b_being_held.html",
			expected:        true,
		},
		{
			name:            "success-abc126_not_being_held",
			inputContestURL: dummyBaseURL + "/contests/abc126",
			mockRequestPath: "/contests/abc126",
			mockStatusCode:  http.StatusOK,
			mockHTMLFile:    "abc126_not_being_held.html",
			expected:        false,
		},
		{
			name:            "failure-xxx999_not_exist",
			inputContestURL: dummyBaseURL + "/contests/xxx999",
			mockRequestPath: "/contests/xxx999",
			mockStatusCode:  http.StatusNotFound,
			mockHTMLFile:    "xxx999_not_exist.html",
			expectedErrMsg:  "could not get HTML",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			html, err := ioutil.ReadFile(path.Join("testdata", "contest", test.mockHTMLFile))
			if err != nil {
				t.Fatal(err)
			}

			defer gock.Off()
			gock.New(dummyBaseURL).
				Get(test.mockRequestPath).
				Reply(test.mockStatusCode).
				AddHeader("Content-Type", "text/html").
				BodyString(string(html))

			c := &Client{baseURL: dummyBaseURL, collector: colly.NewCollector()}
			actual, err := c.IsContestBeingHeld(test.inputContestURL)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err)
				}
				if actual != test.expected {
					t.Fatalf("beingHeld wrong. want=%t, got=%t", test.expected, actual)
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

func TestClient_GetProblemURL(t *testing.T) {
	tests := []struct {
		name string

		inputContest string
		inputProblem string

		mockRequestPath string
		mockStatusCode  int
		mockHTMLFile    string

		expectedProblemURL string
		expectedErrMsg     string
	}{
		{
			name:               "success-abc001",
			inputContest:       "abc001",
			inputProblem:       "C",
			mockRequestPath:    "/contests/abc001/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "abc001.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc001/tasks/abc001_3",
		},
		{
			name:               "success-abc001_case_mess",
			inputContest:       "aBC001",
			inputProblem:       "a",
			mockRequestPath:    "/contests/abc001/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "abc001.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc001/tasks/abc001_1",
		},
		{
			name:               "success-abc124",
			inputContest:       "abc124",
			inputProblem:       "B",
			mockRequestPath:    "/contests/abc124/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "abc124.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/abc124/tasks/abc124_b",
		},
		{
			name:               "success-agc032",
			inputContest:       "agc032",
			inputProblem:       "F",
			mockRequestPath:    "/contests/agc032/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "agc032.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/agc032/tasks/agc032_f",
		},
		{
			name:               "success-arc103",
			inputContest:       "arc103",
			inputProblem:       "E",
			mockRequestPath:    "/contests/arc103/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "arc103.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/arc103/tasks/arc103_c",
		},
		{
			name:               "success-tenka1_2019",
			inputContest:       "tenka1-2019",
			inputProblem:       "E",
			mockRequestPath:    "/contests/tenka1-2019/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "tenka1-2019.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/tenka1-2019/tasks/tenka1_2019_e",
		},
		{
			name:               "success-tenka1_2019",
			inputContest:       "tenka1-2019",
			inputProblem:       "E",
			mockRequestPath:    "/contests/tenka1-2019/tasks",
			mockStatusCode:     http.StatusOK,
			mockHTMLFile:       "tenka1-2019.html",
			expectedProblemURL: "https://dummyatcoder.jp/contests/tenka1-2019/tasks/tenka1_2019_e",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			html, err := ioutil.ReadFile(path.Join("testdata", "problem_list", test.mockHTMLFile))
			if err != nil {
				t.Fatal(err)
			}

			defer gock.Off()
			gock.New(dummyBaseURL).
				Get(test.mockRequestPath).
				Reply(test.mockStatusCode).
				AddHeader("Content-Type", "text/html").
				BodyString(string(html))

			c := &Client{baseURL: dummyBaseURL, collector: colly.NewCollector()}
			problemURL, err := c.GetProblemURL(test.inputContest, test.inputProblem)
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

func TestClient_GetSamples(t *testing.T) {
	tests := []struct {
		name string

		inputProblemURL   string
		inputUseCache     bool
		inputCacheDirPath string

		mockRequestPath string
		mockStatusCode  int
		mockHTMLFile    string

		expectedSamples []Sample
		expectedErrMsg  string
	}{
		{
			name: "success-disable_cache",

			inputProblemURL:   dummyBaseURL + "/contests/abc124/tasks/abc124_b",
			inputUseCache:     false,
			inputCacheDirPath: dummyCacheDirPath,

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
			name: "success-use_cache",

			inputProblemURL:   dummyBaseURL + "/contests/abc124/tasks/abc124_b",
			inputUseCache:     true,
			inputCacheDirPath: dummyCacheDirPath,

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
			name: "success-old_DOM_structure",

			inputProblemURL:   dummyBaseURL + "/contests/abc002/tasks/abc002_c",
			inputUseCache:     false,
			inputCacheDirPath: dummyCacheDirPath,

			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/abc002/tasks/abc002_c",
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
			name: "success-only_one_sample",

			inputProblemURL:   dummyBaseURL + "/contests/kupc2015/tasks/kupc2015_a",
			inputUseCache:     false,
			inputCacheDirPath: dummyCacheDirPath,

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
			name: "failure-nonexistent_problem",

			inputProblemURL:   dummyBaseURL + "/contests/xxx999/tasks/xxx999_x",
			inputUseCache:     false,
			inputCacheDirPath: dummyCacheDirPath,

			mockStatusCode:  http.StatusNotFound,
			mockRequestPath: "contests/xxx999/tasks/xxx999_x",
			mockHTMLFile:    "xxx999x.html",

			expectedErrMsg: "could not get HTML",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := os.RemoveAll(dummyCacheDirPath); err != nil {
					t.Fatalf("failed to remove dummy cache dir: %s", err.Error())
				}
			}()

			html, err := ioutil.ReadFile(path.Join("testdata", "problem", test.mockHTMLFile))
			if err != nil {
				t.Fatal(err)
			}

			defer gock.Off()
			gock.New(dummyBaseURL).
				Get(test.mockRequestPath).
				Reply(test.mockStatusCode).
				AddHeader("Content-Type", "text/html").
				BodyString(string(html))

			if test.inputUseCache {
				if err := os.MkdirAll(dummyCacheDirPath, 0777); err != nil {
					t.Fatalf("failed to create dummy cache dir: %s", err.Error())
				}
				b, err := json.Marshal(test.expectedSamples)
				if err != nil {
					t.Fatalf("failed to marshal samples: %s", err.Error())
				}
				escapedURL := strings.Replace(test.inputProblemURL, "/", "_", -1)
				filename := fmt.Sprintf("%s.json", escapedURL)
				if err := ioutil.WriteFile(path.Join(dummyCacheDirPath, filename), b, 0644); err != nil {
					t.Fatalf("failed to create cache file: %s", err.Error())
				}
			}

			var errBuff bytes.Buffer
			c := &Client{baseURL: dummyBaseURL, collector: colly.NewCollector(), useCache: test.inputUseCache, cacheDirPath: test.inputCacheDirPath, errStream: &errBuff}
			samples, err := c.GetSamples(test.inputProblemURL)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("err should be nil. got: %s", err.Error())
				}
				if errBuff.String() != "" {
					t.Fatalf("errStream should be empty. got: %s", errBuff.String())
				}
				if len(samples) != len(test.expectedSamples) {
					t.Fatalf("length of samples wrong. want=%d, got=%d", len(test.expectedSamples), len(samples))
				}
				for i, expected := range test.expectedSamples {
					actual := samples[i]
					if actual != expected {
						t.Fatalf("%d-th sample wrong.\nwant:\n%+v\ngot:\n%+v", i, expected, actual)
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
				"入力例1": "1 3 5\n",
				"出力例1": "9\n",
				"入力例2": "2 4\n",
				"出力例2": "6\n",
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
				"入力例1": "1 3 5\n",
				"出力例1": "9\n",
				"入力例2": "2 4\n",
			},
			expectedErrMsg: "number of sample elements should be even",
		},
		{
			name: "failure-index_of_入力例_wrong",
			inputElements: map[string]string{
				"入力例2": "1 3 5\n",
				"出力例1": "9\n",
			},
			expectedErrMsg: "could not find '入力例1'",
		},
		{
			name: "failure-index_of_出力例_wrong",
			inputElements: map[string]string{
				"入力例1": "1 3 5\n",
				"出力例2": "9\n",
			},
			expectedErrMsg: "could not find '出力例1'",
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
				"入力例1": strings.Join([]string{
					"4",
					"6 5 6 8",
					"",
				}, "\n"),
				"出力例1": "3\n",
				"入力例2": strings.Join([]string{
					"5",
					"4 5 3 5 4",
					"",
				}, "\n"),
				"出力例2": "3\n",
				"入力例3": strings.Join([]string{
					"5",
					"9 5 6 8 4",
					"",
				}, "\n"),
				"出力例3": "1\n",
			},
		},
		{
			name:            "success-old_DOM_structure",
			inputProblemURL: dummyBaseURL + "/contests/abc002/tasks/abc002_2",
			mockStatusCode:  http.StatusOK,
			mockRequestPath: "contests/abc002/tasks/abc002_2",
			mockHTMLFile:    "abc002c.html",
			expectedSampleElements: map[string]string{
				"入力例1": strings.Join([]string{
					"1 0 3 0 2 5",
					"",
				}, "\n"),
				"出力例1": "5.0\n",
				"入力例2": strings.Join([]string{
					"-1 -2 3 4 5 6",
					"",
				}, "\n"),
				"出力例2": "2.0\n",
				"入力例3": strings.Join([]string{
					"298 520 903 520 4 663",
					"",
				}, "\n"),
				"出力例3": "43257.5\n",
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
			html, err := ioutil.ReadFile(path.Join("testdata", "problem", test.mockHTMLFile))
			if err != nil {
				t.Fatal(err)
			}

			defer gock.Off()
			gock.New(dummyBaseURL).
				Get(test.mockRequestPath).
				Reply(test.mockStatusCode).
				AddHeader("Content-Type", "text/html").
				BodyString(string(html))

			c := &Client{collector: colly.NewCollector()}
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
