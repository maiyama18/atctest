package atcoder

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os/exec"
	"strings"
)

const baseURL = "https://atcoder.jp/contests"

type Sample struct {
	Input  string
	Output string
}

func GetSamples(contest string, problem string) ([]Sample, error) {
	url := getURL(contest, problem)
	elements, err := fetchSampleElements(url)
	if err != nil {
		return nil, err
	}

	return constructSamples(elements)
}

func FormatSamples(samples []Sample) string {
	var out bytes.Buffer
	for i, sample := range samples {
		out.WriteString(fmt.Sprintf("# sample %d\n", i+1))
		out.WriteString("input:\n")
		out.WriteString(sample.Input)
		out.WriteString("output:\n")
		out.WriteString(sample.Output)
		out.WriteString("\n")
	}

	return out.String()
}

func Check(samples []Sample, command string) {
	fields := strings.Fields(command)
	c := fields[0]
	args := fields[1:]

	for i, sample := range samples {
		cmd := exec.Command(c, args...)
		cmd.Stdin = strings.NewReader(sample.Input)

		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		result := "x"
		if string(out) == sample.Output {
			result = "o"
		}
		fmt.Printf("sample %d -> %s\n", i+1, result)
		fmt.Print("expected:", sample.Output)
		fmt.Print("actual:", string(out))
		fmt.Println()
	}
}

func getURL(contest string, problem string) string {
	return fmt.Sprintf("%s/%s/tasks/%s_%s", baseURL, strings.ToLower(contest), strings.ToLower(contest), strings.ToLower(problem))
}

func fetchSampleElements(url string) (map[string]string, error) {
	c := colly.NewCollector()

	elements := make(map[string]string)
	c.OnHTML(`pre`, func(e *colly.HTMLElement) {
		title := e.DOM.Parent().Find("h3").Text()
		if strings.HasPrefix(title, "入力例") || strings.HasPrefix(title, "出力例") {
			elements[title] = e.Text
		}
	})

	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("could not get HTML: %s", url)
	}

	return elements, nil
}

func constructSamples(elements map[string]string) ([]Sample, error) {
	if len(elements) == 0 {
		return nil, errors.New("no sample elements found")
	}
	if len(elements)%2 != 0 {
		return nil, fmt.Errorf("number of sample elements should be even because it consists of pair of input/output. got: %d", len(elements))
	}

	numSamples := len(elements) / 2
	samples := make([]Sample, numSamples)
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
