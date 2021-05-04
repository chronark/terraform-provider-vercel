package alias

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type CreateOrUpdateAlias struct {
	Domain   string `json:"domain"`
	Redirect string `json:"redirect"`
}

type Alias struct {
	CreatedAt           int64  `json:"createdAt"`
	Domain              string `json:"domain"`
	Target              string `json:"target"`
	ConfiguredBy        string `json:"configuredBy"`
	ConfiguredChangedAt int64  `json:"configuredChangedAt"`
}

type Handler struct {
	Api httpApi.API
}

func (h *Handler) Create(projectId string, alias CreateOrUpdateAlias, teamId string) error {
	url := fmt.Sprintf("/v1/projects/%s/alias", projectId)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request(http.MethodPost, url, alias)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var createdAliases []Alias
	err = json.NewDecoder(res.Body).Decode(&createdAliases)
	if err != nil {
		return err
	}
	log.Printf("%+v\n\n", createdAliases)

	return nil
}

func (h *Handler) Update(projectId string, alias CreateOrUpdateAlias, teamId string) error {
	url := fmt.Sprintf("/v1/projects/%s/alias", projectId)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := h.Api.Request("PATCH", url, alias)
	if err != nil {
		return fmt.Errorf("Unable to update env: %w", err)
	}
	defer res.Body.Close()
	return nil
}
func (h *Handler) Delete(projectId, envKey string, teamId string) error {
	url := fmt.Sprintf("/v1/projects/%s/alias?domain", projectId)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete domain: %w", err)
	}
	defer res.Body.Close()
	return nil
}
