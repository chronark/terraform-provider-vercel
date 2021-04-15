package team

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type Team struct {
	Id        string    `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	CreatorId string    `json:"creatorId"`
	Created   time.Time `json:"created"`
	Avatar    string    `json:"avatar"`
}

type Handler struct {
	Api httpApi.API
}

func (h *Handler) Read(slug string) (user Team, err error) {
	res, err := h.Api.Request("GET", fmt.Sprintf("/v1/teams/?slug=%s", slug), nil)
	if err != nil {
		return Team{}, fmt.Errorf("Unable to fetch team from vercel: %w", err)
	}
	defer res.Body.Close()

	var team Team
	err = json.NewDecoder(res.Body).Decode(&team)
	if err != nil {
		return Team{}, fmt.Errorf("Unable to unmarshal team: %w", err)
	}
	return team, nil
}
