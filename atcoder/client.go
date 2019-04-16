package atcoder

import (
	"errors"
	"fmt"
	"net/http"
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

	problemURL string
}

func NewClient(baseURL, contest, problem string) (*Client, error) {
	c := &Client{baseURL: baseURL, contest: contest, problem: problem}
	url, err := c.setProblemURL()
	if err != nil {
		return nil, err
	}
	c.problemURL = url

	return c, nil
}

func (c *Client) GetSamples() ([]Sample, error) {
	elements, err := c.fetchSampleElements()
	if err != nil {
		return nil, err
	}

	return c.constructSamples(elements)
}

func (c *Client) setProblemURL() (string, error) {
	u1 := c.urlType1()
	if isValidURL(u1) {
		return u1, nil
	}

	u2 := c.urlType2()
	if isValidURL(u2) {
		return u2, nil
	}

	return "", fmt.Errorf("ERROR: could not find problem page for problem '%s' of contest '%s'", c.problem, c.contest)
}

func (c *Client) urlType1() string {
	contestStr := strings.ToLower(c.contest)
	problemStr := strings.ToLower(c.problem)
	return fmt.Sprintf("%s/contests/%s/tasks/%s_%s", c.baseURL, contestStr, contestStr, problemStr)
}

func (c *Client) urlType2() string {
	contestStr := strings.ToLower(c.contest)

	problemStr, ok := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5",
		"f": "6",
	}[strings.ToLower(c.problem)]
	if !ok {
		problemStr = "0"
	}

	return fmt.Sprintf("%s/contests/%s/tasks/%s_%s", c.baseURL, contestStr, contestStr, problemStr)
}

func (c *Client) fetchSampleElements() (map[string]string, error) {
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

	if err := cl.Visit(c.problemURL); err != nil {
		return nil, fmt.Errorf("ERROR: could not get HTML: %s", c.problemURL)
	}

	return elements, nil
}

func (c *Client) constructSamples(elements map[string]string) ([]Sample, error) {
	if len(elements) == 0 {
		return nil, errors.New("ERROR: no sample elements found")
	}
	if len(elements)%2 != 0 {
		return nil, fmt.Errorf("ERROR: number of sample elements should be even because it consists of pair of input/output. got: %d", len(elements))
	}

	numSamples := len(elements) / 2
	samples := make([]Sample, numSamples)
	for i := 1; i <= numSamples; i++ {
		inputKey := fmt.Sprintf("入力例 %d", i)
		outputKey := fmt.Sprintf("出力例 %d", i)

		input, ok := elements[inputKey]
		if !ok {
			return nil, fmt.Errorf("ERROR: could not find '%s' in HTML", inputKey)
		}
		output, ok := elements[outputKey]
		if !ok {
			return nil, fmt.Errorf("ERROR: could not find '%s' in HTML", outputKey)
		}

		samples[i-1] = Sample{Input: input, Output: output}
	}

	return samples, nil
}

func isValidURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}
