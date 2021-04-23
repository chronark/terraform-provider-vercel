package httpApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API interface {
	Request(method string, path string, body interface{}) (*http.Response, error)
}

type Api struct {
	url        string
	httpClient *http.Client
	userAgent  string
	token      string
}

func New(token string) API {
	return &Api{
		url:        "https://api.vercel.com",
		httpClient: &http.Client{},
		userAgent:  "chronark/terraform-provider-vercel",
		token:      token,
	}
}

// https://vercel.com/docs/api#api-basics/errors
type VercelError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Api) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
}

func (c *Api) do(req *http.Request) (*http.Response, error) {
	c.setHeaders(req)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to perform request: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		defer res.Body.Close()

		var x map[string]interface{}
		_ = json.NewDecoder(res.Body).Decode(&x)

		// var vercelError VercelError
		// err = json.NewDecoder(res.Body).Decode(&vercelError)
		// if err != nil {
		// 	return nil, fmt.Errorf("Error during http request: Unable to extract error: %w", err)

		// }
		return nil, fmt.Errorf("Error during http request: %+v", x)
	}
	return res, nil
}

func (c *Api) Request(method string, path string, body interface{}) (*http.Response, error) {
	var payload io.Reader = nil
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.url, path), payload)
	if err != nil {
		return nil, err
	}
	res, err := c.do(req)

	if err != nil {
		return &http.Response{}, fmt.Errorf("Unable to request resource: [%s] %s with payload {%+v}: %w", method, path, payload, err)
	}
	return res, nil

}
