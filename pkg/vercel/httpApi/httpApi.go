package httpApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Api struct {
	url        string
	httpClient *http.Client
	userAgent  string
	token      string
}

func New(token string) *Api {
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

func (c *Api) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.url, path), nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)

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
		log.Printf("%+v\n", x)

		var vercelError VercelError
		err = json.NewDecoder(res.Body).Decode(&vercelError)
		if err != nil {
			var body map[string]interface{}
			_ = json.NewDecoder(res.Body).Decode(&body)

			return nil, fmt.Errorf("Request was not successfull [ %s ], but I could not unmarshal the response body: %w. Request was: %+v. Raw response body was: %s", res.Status, err, req, body)
		}
		return nil, fmt.Errorf("Error during http request: [ %s ] - %s", vercelError.Error.Code, vercelError.Error.Message)
	}
	return res, nil
}

func (c *Api) Post(path string, body interface{}) (*http.Response, error) {
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.url, path), bytes.NewBuffer(encodedBody))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Api) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s", c.url, path), nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
func (c *Api) Patch(path string, body interface{}) (*http.Response, error) {
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s", c.url, path), bytes.NewBuffer(encodedBody))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
