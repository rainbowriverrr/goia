package requests

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CompletionRequest struct {
	// URL is the URL of the OpenAI API endpoint.
	URL string `json:"-"`
	//Bearer is the bearer token to use for authentication.
	Bearer string `json:"-"`
	// Model is the model to use for completion.
	Model string `json:"model"`
	// Prompt is the prompt to use for completion.
	Prompt string `json:"prompt"`
	// Stop is a list of strings that the API will use to stop generating further text.
	Stop []string `json:"stop,omitempty"`
	// MaxTokens is the maximum number of tokens to generate. (one token = 4 characters ish)
	MaxTokens int `json:"max_tokens"`
	// Completions is the number of completions to generate.
	Completions int `json:"n"`
	// Temperature is a value between 0 and 1 that controls randomness. 1 is the most random.
	Temperature float64 `json:"temperature"`
	// TopP is the cumulative probability threshold for top-k sampling. 1 is the least conservative.
	TopP float64 `json:"top_p"`
	//LogProbs is the number of highest probability tokens to return. If null, no log-probs are returned. The default is 0, the maximum is 5.
	LogProbs int `json:"logprobs"`
	// Echo is whether to return the prompt in addition to the completion. If null, defaults to false.
	Echo bool `json:"echo"`
	//BestOf is the number of completions the model will generate before returning the best one. If null, defaults to 1.
	BestOf int `json:"best_of"`
	// FrequencyPenalty is a value between 0 and 1 that penalizes new tokens based on their existing frequency in the text so far.  Higher = less repetitive.
	FrequencyPenalty float64 `json:"frequency_penalty"`
	//Presence Penalty is a value between -2.0 to 2.0 that penalizes new tokens based on whether they appear in the text so far. Higher = more likely to talk about new topics.
	PresencePenalty float64 `json:"presence_penalty"`
	// LogitBias is a dictionary of token strings to bias values, representing how likely the model is to output the given tokens. -100 to 100. TODO: implement this use tokenizer tool.
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// User is the ID of the user who is using the API. Good for logging and moderation.
	User string `json:"user"`
	//Stream is whether to stream back partial progress. If null, defaults to false. Terminated by data: [DONE].
	Stream bool `json:"stream"`
}

func NewCompletionRequest() *CompletionRequest {
	toReturn := &CompletionRequest{}
	toReturn.MakeDefault()
	return toReturn
}

func (c *CompletionRequest) SetBearer(token string) {
	c.Bearer = token
}

// SetModel sets the model to use for completion.
func (c *CompletionRequest) SetModel(model string) {
	c.Model = model
}

// SetPrompt sets the prompt for the completion request.
func (c *CompletionRequest) SetPrompt(prompt string) {
	c.Prompt = prompt
}

// AddStop adds a stop string to the completion request.  This will tell the model when to stop generating.
func (c *CompletionRequest) AddStop(stop string) {
	if c.Stop == nil {
		c.Stop = make([]string, 0)
	}
	c.Stop = append(c.Stop, stop)
}

// SetMaxTokens sets the maximum number of tokens to generate.  One token is about 4 characters.
func (c *CompletionRequest) SetMaxTokens(maxTokens int) {
	c.MaxTokens = maxTokens
}

// SetCompletions sets the number of completions to generate.
func (c *CompletionRequest) SetCompletions(completions int) {
	c.Completions = completions
}

// SetTemperature sets the temperature of the model.  1 is the most random.
func (c *CompletionRequest) SetTemperature(temperature float64) {
	if temperature < 0 {
		temperature = 0
	} else if temperature > 1 {
		temperature = 1
	}
	c.Temperature = temperature
}

// SetTopP sets the cumulative probability threshold for top-k sampling.  1 is the least conservative.
func (c *CompletionRequest) SetTopP(topP float64) {
	if topP < 0 {
		topP = 0
	} else if topP > 1 {
		topP = 1
	}
	c.TopP = topP
}

// SetLogProbs sets the number of highest probability tokens to return. If null, no log-probs are returned. The default is 0, the maximum is 5.
func (c *CompletionRequest) SetLogProbs(logProbs int) {
	if logProbs < 0 {
		logProbs = 0
	} else if logProbs > 5 {
		logProbs = 5
	}
	c.LogProbs = logProbs
}

// SetEcho sets whether to return the prompt in addition to the completion. If null, defaults to false.
func (c *CompletionRequest) SetEcho(echo bool) {
	c.Echo = echo
}

// SetBestOf sets the number of completions the model will generate before returning the best one. If null, defaults to 1.
func (c *CompletionRequest) SetBestOf(bestOf int) {
	if bestOf < 1 {
		bestOf = 1
	}
	c.BestOf = bestOf
}

// SetFrequencyPenalty sets the frequency penalty, between 0 - 1
func (c *CompletionRequest) SetFrequencyPenalty(freqp float64) {
	if freqp < 0 {
		freqp = 0
	} else if freqp > 1 {
		freqp = 1
	}
	c.FrequencyPenalty = freqp
}

// SetPresencePenalty sets the presence penalty of the request.  This penalizes new tokens based on whether they have already appeared in the text.
func (c *CompletionRequest) SetPresencePenalty(pp float64) {
	if pp < -2 {
		pp = -2
	} else if pp > 2 {
		pp = 2
	}
	c.PresencePenalty = pp
}

// AddBias adds a bias to a given token. This will make it more likely to appear.
// If the token already exists in the bias, the bias will be added to the existing bias.
func (c *CompletionRequest) AddBias(tok string, bias int) {

	if c.LogitBias == nil {
		c.LogitBias = make(map[string]int)
	}

	currbias, exists := c.LogitBias[tok]
	if !exists {
		//checks the bias to make sure it's in the range of -100 to 100
		//sets it to be within range if it currently isnt
		if bias > 100 {
			bias = 100
		} else if bias < -100 {
			bias = -100
		}
		c.LogitBias[tok] = bias
	} else {
		//checks if the bias will make the bias out of range
		//if it does, it sets it to the max or min value
		if currbias+bias > 100 {
			c.LogitBias[tok] = 100
		} else if currbias+bias < -100 {
			c.LogitBias[tok] = -100
		} else {
			c.LogitBias[tok] += bias
		}
	}
}

// RemoveBias removes a bias from a given token. This will make it less likely to appear.
// If the token already exists in the bias, the bias will be subtracted from the existing bias.
func (c *CompletionRequest) RemoveBias(tok string, bias int) {

	if c.LogitBias == nil {
		c.LogitBias = make(map[string]int)
	}

	currbias, exists := c.LogitBias[tok]
	if !exists {
		//checks the bias to make sure it's in the range of -100 to 100
		//sets it to be within range if it currently isnt
		if bias > 100 {
			bias = 100
		} else if bias < -100 {
			bias = -100
		}
		c.LogitBias[tok] = bias
	} else {
		//checks if the bias will make the bias out of range
		//if it does, it sets it to the max or min value
		if currbias-bias > 100 {
			c.LogitBias[tok] = 100
		} else if currbias-bias < -100 {
			c.LogitBias[tok] = -100
		} else {
			c.LogitBias[tok] -= bias
		}
	}
}

// SetUser sets the user for the completion request. This can be any string, use it for tracking and moderation.
func (c *CompletionRequest) SetUser(user string) {
	c.User = user
}

// SetStream sets whether to stream the response. If null, defaults to false.
func (c *CompletionRequest) SetStream(stream bool) {
	c.Stream = stream
}

// MakeDefault makes the completion request into a default request.
func (c *CompletionRequest) MakeDefault() {
	c.URL = "https://api.openai.com/v1/completions"
	c.Model = "text-davinci-003"
	c.Prompt = ""
	c.Stop = []string{"User:"}
	c.MaxTokens = 16
	c.Completions = 1
	c.Temperature = 1
	c.TopP = 1
	c.LogProbs = 0
	c.Echo = false
	c.BestOf = 1
	c.FrequencyPenalty = 0
	c.PresencePenalty = 0
	c.LogitBias = nil
	c.Stream = false
	c.User = ""
}

// GetRequest converts the completionrequest into an http.Request
func (c *CompletionRequest) GetRequest() (*http.Request, error) {
	json, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Bearer)
	return req, nil
}
