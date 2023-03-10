package triviadb

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type (
	Result struct {
		Category         string   `json:"category,omitempty"`
		Type             string   `json:"type"`
		Difficulty       string   `json:"difficulty"`
		Question         string   `json:"question"`
		CorrectAnswer    string   `json:"correct_answer"`
		IncorrectAnswers []string `json:"incorrect_answers"`
		CustomerAnswer   string   `json:"canswer,omitempty"`
	}

	Results struct {
		ResponseCode int      `json:"response_code"`
		Results      []Result `json:"results"`
	}
)

func readResults(r io.ReadCloser) ([]Result, error) {
	res := &Results{}
	f, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, res)
	if err != nil {
		return nil, err
	}
	return res.Results, nil
}

func GetResults(url string) ([]Result, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	c, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", url, err)
	}
	defer c.Body.Close()
	return readResults(c.Body)
}

func (r *Result) GetAnswers() []string {
	rand.Seed(time.Now().UnixNano())
	answers := []string{}
	answers = append(answers, html.UnescapeString(r.CorrectAnswer))
	for _, ans := range r.IncorrectAnswers {
		answers = append(answers, html.UnescapeString(ans))
	}
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	return answers
}

func (r *Result) GetCorrectAnswer() string {
	return html.UnescapeString(r.CorrectAnswer)
}
