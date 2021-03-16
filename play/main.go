package main

import (
	"encoding/json"
	"log"

	"github.com/chronark/terraform-provider-vercel/internal/vercel"
)

type User struct {
	User struct {
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
	} `json:"user"`
}

func main() {
	client := vercel.New("wsByP9ptGqn7snGvvY00aDzn")

	res, err := client.Call("GET", "/www/user")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	log.Println()
	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	log.Println(user)
}
