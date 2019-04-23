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

//var dummyCacheDirPath = path.Join("testdata", "cache")
//
//var dummySamples = []Sample{
//	{Input: "input1\n", Output: "output1\n"},
//	{Input: "input2\n", Output: "output2\n"},
//}

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

			c := &Client{baseURL: dummyBaseURL}
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

//func TestClient_getCachedSamples(t *testing.T) {
//	tests := []struct {
//		name string
//
//		inputContest      string
//		inputProblem      string
//		inputCacheDirPath string
//
//		expectedSamples []Sample
//		expectedSuccess bool
//	}{
//		{
//			name:              "success",
//			inputContest:      "abc120",
//			inputProblem:      "c",
//			inputCacheDirPath: dummyCacheDirPath,
//			expectedSamples:   dummySamples,
//			expectedSuccess:   true,
//		},
//		{
//			name:              "failed-cache_dir_path_not_exist",
//			inputContest:      "abc120",
//			inputProblem:      "c",
//			inputCacheDirPath: "/nonexistent",
//			expectedSuccess:   false,
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			if err := os.MkdirAll(dummyCacheDirPath, 0777); err != nil {
//				t.Fatalf("failed to create dummy cache dir: %s", err.Error())
//			}
//			b, err := json.Marshal(test.expectedSamples)
//			if err != nil {
//				t.Fatalf("failed to marshal samples: %s", err.Error())
//			}
//			filename := fmt.Sprintf("%s-%s.json", test.inputContest, test.inputProblem)
//			if err := ioutil.WriteFile(path.Join(dummyCacheDirPath, filename), b, 0644); err != nil {
//				t.Fatalf("failed to create cache file: %s", err.Error())
//			}
//
//			defer func() {
//				if err := os.RemoveAll(dummyCacheDirPath); err != nil {
//					t.Fatalf("failed to remove dummy cache dir after test: %s", dummyCacheDirPath)
//				}
//			}()
//
//			c := &Client{cacheDirPath: test.inputCacheDirPath}
//			samples, ok := c.getCachedSamples()
//			if test.expectedSuccess {
//				if !ok {
//					t.Fatalf("should success to get cache, but failed")
//				}
//
//				if len(samples) != len(test.expectedSamples) {
//					t.Fatalf("length of samples wrong. want=%d, got=%d", len(test.expectedSamples), len(samples))
//				}
//				for i, expected := range test.expectedSamples {
//					actual := samples[i]
//					if actual != expected {
//						t.Fatalf("%d-th sample wrong. want=%+v, got=%+v", i, expected, actual)
//					}
//				}
//			} else {
//				if ok {
//					t.Fatalf("should fail to get cache, but succeeded")
//				}
//			}
//		})
//	}
//}

//func TestClient_cacheSamples(t *testing.T) {
//	tests := []struct {
//		name string
//
//		inputContest      string
//		inputProblem      string
//		inputCacheDirPath string
//		inputSamples      []Sample
//
//		expectedErrMsg string
//	}{
//		{
//			name:              "success",
//			inputContest:      "abc120",
//			inputProblem:      "c",
//			inputCacheDirPath: dummyCacheDirPath,
//			inputSamples: []Sample{
//				{Input: "input1\n", Output: "output1\n"},
//				{Input: "input2\n", Output: "output2\n"},
//			},
//		},
//		{
//			name:              "success-cache_dir_not_exist",
//			inputContest:      "abc120",
//			inputProblem:      "c",
//			inputCacheDirPath: path.Join("testdata", "new_cache_dir"),
//			inputSamples:      dummySamples,
//		},
//		{
//			name:              "failed-cache_dir_failed_no_permission",
//			inputContest:      "abc120",
//			inputProblem:      "c",
//			inputCacheDirPath: path.Join("/sys", "new_cache_dir"),
//			inputSamples:      dummySamples,
//			expectedErrMsg:    "permission denied",
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			if err := os.MkdirAll(dummyCacheDirPath, 0777); err != nil {
//				t.Fatal("failed to create dummy cache dir")
//			}
//			defer func() {
//				if err := os.RemoveAll(dummyCacheDirPath); err != nil {
//					t.Fatalf("failed to remove dummy cache dir after test: %s", dummyCacheDirPath)
//				}
//				if test.inputCacheDirPath != dummyCacheDirPath {
//					if err := os.RemoveAll(test.inputCacheDirPath); err != nil {
//						t.Fatalf("failed to remove cache dir after test: %s", test.inputCacheDirPath)
//					}
//				}
//			}()
//
//			c := &Client{cacheDirPath: test.inputCacheDirPath}
//			err := c.cacheSamples(test.inputSamples)
//			if test.expectedErrMsg == "" {
//				if err != nil {
//					t.Fatalf("err should be nil. got: %s", err)
//				}
//
//				filename := fmt.Sprintf("%s-%s.json", c.contest, c.problem)
//				b, err := ioutil.ReadFile(path.Join(c.cacheDirPath, filename))
//				if err != nil {
//					t.Fatalf("failed to read cached file")
//				}
//				var samples []Sample
//				if err := json.Unmarshal(b, &samples); err != nil {
//					t.Fatalf("failed to unmarshal samples")
//				}
//
//				if len(samples) != len(test.inputSamples) {
//					t.Fatalf("length of samples wrong. want=%d, got=%d", len(test.inputSamples), len(samples))
//				}
//				for i, expected := range test.inputSamples {
//					actual := samples[i]
//					if actual != expected {
//						t.Fatalf("%d-th sample wrong. want=%+v, got=%+v", i, expected, actual)
//					}
//				}
//			} else {
//				if err == nil {
//					t.Fatal("err should not be nil. got: nil")
//				}
//				if !strings.Contains(err.Error(), test.expectedErrMsg) {
//					t.Fatalf("expect '%s' to contain '%s'", err.Error(), test.expectedErrMsg)
//				}
//			}
//		})
//	}
//}
