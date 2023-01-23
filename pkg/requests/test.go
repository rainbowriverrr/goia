package requests

import (
	"encoding/json"
	"net/http"
	"strings"
)

type TestRequest struct {
	URL              string   `json:"-"`
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	Stop             []string `json:"stop"`
	MaxTokens        int      `json:"max_tokens"`
	Temperature      float64  `json:"temperature"`
	TopP             float64  `json:"top_p"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	PresencePenalty  float64  `json:"presence_penalty"`
}

func NewTestRequest() *TestRequest {
	return &TestRequest{
		URL:              "https://api.openai.com/v1/completions",
		Model:            "text-davinci-003",
		Prompt:           "You: What have you been up to?\nFriend: Watching old movies.\nYou: Did you watch anything interesting?\nFriend:",
		Temperature:      0.5,
		MaxTokens:        60,
		TopP:             1,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0.0,
		Stop:             []string{"You:"},
	}
}

// GetRequest converts a TestRequest to a http.Request by marshalling the TestRequest into JSON body.
func (t *TestRequest) GetRequest() (*http.Request, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, t.URL, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
