package alias

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type CreateOrUpdateAlias struct {
	Domain   string `json:"domain"`
	Redirect string `json:"redirect"`
}

type Alias struct {
	ConfiguredBy        string `json:"configuredBy"`
	ConfiguredChangedAt int    `json:"configuredChangedAt"`
	GitBranch           string `json:"gitBranch"`
	ProjectId           string `json:"projectId"`
	CreatedAt           int64  `json:"createdAt"`
	Domain              string `json:"domain"`
	Redirect            string `json:"redirect"`
	RedirectStatusCode  int    `json:"redirectStatusCode"`
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
	return nil
}

func (h *Handler) Read(projectId string, domain, teamId string) (Alias, error) {
	url := fmt.Sprintf("/v1/projects/%s", projectId)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := h.Api.Request("GET", url, nil)
	if err != nil {
		return Alias{}, fmt.Errorf("Unable to fetch project from vercel: %w", err)
	}
	defer res.Body.Close()
	type Project struct {
		Aliases []Alias `json:"alias"`
	}
	project := Project{}

	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {

		return Alias{}, fmt.Errorf("Unable to unmarshal response from %s: %w", url, err)
	}

	for _, alias := range project.Aliases {
		if alias.Domain == domain {
			return alias, nil
		}
	}

	return Alias{}, fmt.Errorf("No alias with domain: %s found", domain)
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
func (h *Handler) Delete(projectId, domain string, teamId string) error {
	url := fmt.Sprintf("/v1/projects/%s/alias?domain=%s", projectId, domain)
	if teamId != "" {
		url = fmt.Sprintf("%s&teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete domain: %w", err)
	}
	defer res.Body.Close()
	return nil
}
