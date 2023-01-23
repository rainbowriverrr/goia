package client

import (
	"io"
	"net/http"

	"github.com/rainbowriverrr/goai/pkg/requests"
)

type Client struct {
	bearerToken string
	baseURL     string
	HttpClient  *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		HttpClient:  &http.Client{},
		bearerToken: token,
	}
}

func (c *Client) SetBearerToken(token string) {
	c.bearerToken = token
}

// SendCompletionRequest sends the completion request to its endpoint and returns the response body
// TODO: Implement for my request interface
func (c *Client) SendCompletionRequest(req *requests.CompletionRequest) ([]byte, error) {
	req.SetBearer(c.bearerToken)
	httpRequest, err := req.GetRequest()
	if err != nil {
		println("error transforming request")
		return nil, err
	}
	resp, err := c.HttpClient.Do(httpRequest)
	if err != nil {
		println("error sending request")
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("error reading response body")
		return nil, err
	}
	return body, nil
}

// func (c *Client) TestRequestResponse() (*http.Response, error) {
// 	req, err := requests.NewTestRequest().GetRequest()
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("Authorization", "Bearer "+c.bearerToken)

// 	resp, _ := c.HttpClient.Do(req)

// 	return resp, nil
// }

// func GetResponseBody(resp *http.Response) ([]byte, error) {
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return body, nil
// }
