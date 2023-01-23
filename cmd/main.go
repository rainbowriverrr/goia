package main

import (
	"net/http"

	"github.com/rainbowriverrr/goai/pkg/client"
	"github.com/rainbowriverrr/goai/pkg/requests"
)

func main() {
	c := client.Client{
		HttpClient: &http.Client{},
	}
	c.SetBearerToken("sdfs")
	req := requests.NewCompletionRequest()
	req.SetPrompt("User: Write a quick introduction, introducing River.  River is a software engineer, the one who programmed this API to call you.  They are a young computer science student who is transgender.  Make this introduction sound professional. \n AI:")
	req.SetMaxTokens(60)
	req.SetTemperature(0.9)
	req.SetFrequencyPenalty(0.6)
	body, _ := c.SendCompletionRequest(req)
	println(string(body))
}
