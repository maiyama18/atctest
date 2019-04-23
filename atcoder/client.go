package atcoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

type Sample struct {
	Input  string
	Output string
}

type Client struct {
	baseURL string
	contest string
	problem string

	useCache     bool
	cacheDirPath string

	outStream io.Writer
	errStream io.Writer
}

func NewClient(baseURL string, useCache bool, cacheDirPath string, outStream, errStream io.Writer) *Client {
	return &Client{
		baseURL:      baseURL,
		useCache:     useCache,
		cacheDirPath: cacheDirPath,
		outStream:    outStream,
		errStream:    errStream,
	}
}

func (c *Client) GetProblemURL(contest, problem string) (string, error) {
	lowerContest := strings.ToLower(contest)
	lowerProblem := strings.ToLower(problem)
	u1 := c.urlType1(lowerContest, lowerProblem)
	if isValidURL(u1) {
		return u1, nil
	}

	u2 := c.urlType2(lowerContest, lowerProblem)
	if isValidURL(u2) {
		return u2, nil
	}

	return "", fmt.Errorf("could not find problem page for problem '%s' of contest '%s'", c.problem, c.contest)
}

func (c *Client) GetSamples(problemURL string) ([]Sample, error) {
	if c.useCache {
		if samples, ok := c.getCachedSamples(); ok {
			return samples, nil
		}
	}

	elements, err := c.fetchSampleElements(problemURL)
	if err != nil {
		return nil, err
	}

	samples, err := c.constructSamples(elements)
	if err != nil {
		return nil, err
	}

	if c.useCache {
		if err := c.cacheSamples(samples); err != nil {
			_, _ = io.WriteString(c.errStream, err.Error())
		}
	}

	return samples, nil
}

func (c *Client) getCachedSamples() ([]Sample, bool) {
	_, err := os.Stat(c.cacheDirPath)
	if err != nil {
		return nil, false
	}

	filename := fmt.Sprintf("%s-%s.json", c.contest, c.problem)
	bytes, err := ioutil.ReadFile(path.Join(c.cacheDirPath, filename))
	if err != nil {
		return nil, false
	}

	var samples []Sample
	if err := json.Unmarshal(bytes, &samples); err != nil {
		return nil, false
	}

	return samples, true
}

func (c *Client) cacheSamples(samples []Sample) error {
	_, err := os.Stat(c.cacheDirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(c.cacheDirPath, 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	bytes, err := json.Marshal(samples)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s-%s.json", c.contest, c.problem)

	return ioutil.WriteFile(path.Join(c.cacheDirPath, filename), bytes, 0644)
}

func (c *Client) fetchSampleElements(problemURL string) (map[string]string, error) {
	cl := colly.NewCollector()

	elements := make(map[string]string)
	cl.OnHTML(`pre`, func(e *colly.HTMLElement) {
		title := e.DOM.Parent().Find("h3").Text()
		if strings.HasPrefix(title, "入力例") || strings.HasPrefix(title, "出力例") {
			elements[title] = e.Text
		} else {
			title := e.DOM.Parent().Parent().Find("h3").Text()
			if strings.HasPrefix(title, "入力例") || strings.HasPrefix(title, "出力例") {
				elements[title] = e.Text
			}
		}
	})

	if err := cl.Visit(problemURL); err != nil {
		return nil, fmt.Errorf("could not get HTML: %s", problemURL)
	}

	return elements, nil
}

func (c *Client) constructSamples(elements map[string]string) ([]Sample, error) {
	if len(elements) == 0 {
		return nil, errors.New("no sample elements found")
	}
	if len(elements)%2 != 0 {
		return nil, fmt.Errorf("number of sample elements should be even because it consists of pair of input/output. got: %d", len(elements))
	}

	numSamples := len(elements) / 2
	samples := make([]Sample, numSamples)

	// for html which only has one pair without numbering ["入力例", "出力例"] (without numbering)
	if numSamples == 1 {
		if input, ok := elements["入力例"]; ok {
			if output, ok := elements["出力例"]; ok {
				samples[0] = Sample{Input: input, Output: output}
				return samples, nil
			}
		}
	}

	// for html which has pairs of samples with numbering ["入力例 1", "出力例 1", "入力例 2", ...]
	for i := 1; i <= numSamples; i++ {
		inputKey := fmt.Sprintf("入力例 %d", i)
		outputKey := fmt.Sprintf("出力例 %d", i)

		input, ok := elements[inputKey]
		if !ok {
			return nil, fmt.Errorf("could not find '%s' in HTML", inputKey)
		}
		output, ok := elements[outputKey]
		if !ok {
			return nil, fmt.Errorf("could not find '%s' in HTML", outputKey)
		}

		samples[i-1] = Sample{Input: input, Output: output}
	}

	return samples, nil
}

func (c *Client) urlType1(contest, problem string) string {
	return fmt.Sprintf("%s/contests/%s/tasks/%s_%s", c.baseURL, contest, contest, problem)
}

func (c *Client) urlType2(contest, problem string) string {
	problemStr, ok := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5",
		"f": "6",
	}[problem]
	if !ok {
		problemStr = "0"
	}

	return fmt.Sprintf("%s/contests/%s/tasks/%s_%s", c.baseURL, contest, contest, problemStr)
}

func isValidURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}
