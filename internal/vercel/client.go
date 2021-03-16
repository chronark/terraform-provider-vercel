package vercel

import (
	"encoding/json"
	"fmt"
	"net/http"
)




type Client struct {
	url string
	httpClient *http.Client
	userAgent string
}


func New() *Client {
	return &Client{
		url: "https://api.vercel.com",
		httpClient: &http.Client{},
		userAgent: "chronark/terraform-provider-vercel",
	}
}


// https://vercel.com/docs/api#api-basics/errors
type VercelError struct {
	Error struct {
		Code string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) Call(method string, path string)(*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.url, path), nil)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		var vercelError VercelError
		err = json.NewDecoder(res.Body).Decode(&vercelError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error during http request: [ %s ] - %s", vercelError.Error.Code, vercelError.Error.Message)
	}
	return res, nil

}