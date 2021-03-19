package user

import (
	"encoding/json"
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type User struct {
	UID             string `json:"uid"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Avatar          string `json:"avatar"`
	PlatformVersion int    `json:"platformVersion"`
	Billing         struct {
		Plan        string      `json:"plan"`
		Period      interface{} `json:"period"`
		Trial       interface{} `json:"trial"`
		Cancelation interface{} `json:"cancelation"`
		Addons      interface{} `json:"addons"`
	} `json:"billing"`
	Bio      string `json:"bio"`
	Website  string `json:"website"`
	Profiles []struct {
		Service string `json:"service"`
		Link    string `json:"link"`
	} `json:"profiles"`
}

type UserHandler struct {
	Api httpApi.API
}

func (p *UserHandler) Read() (user User, err error) {
	res, err := p.Api.Request("GET", "/www/user", nil)
	if err != nil {
		return User{}, fmt.Errorf("Unable to fetch user from vercel: %w", err)
	}
	defer res.Body.Close()

	type UserResponse struct {
		User User `json:"user"`
	}

	var userResponse UserResponse
	err = json.NewDecoder(res.Body).Decode(&userResponse)
	if err != nil {
		return User{}, fmt.Errorf("Unable to unmarshal project: %w", err)
	}
	return userResponse.User, nil
}
