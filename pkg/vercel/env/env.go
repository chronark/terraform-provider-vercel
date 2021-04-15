package env

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type CreateOrUpdateEnv struct {
	// The type can be `plain`, `secret`, or `system`.
	Type string `json:"type"`

	// The name of the environment variable.
	Key string `json:"key"`

	// If the type is `plain`, a string representing the value of the environment variable.
	// If the type is `secret`, the secret ID of the secret attached to the environment variable.
	// If the type is `system`, the name of the System Environment Variable.
	Value string `json:"value"`

	// 	The target can be a list of `development`, `preview`, or `production`.
	Target []string `json:"target"`
}

type Env struct {
	Type            string      `json:"type"`
	ID              string      `json:"id"`
	Key             string      `json:"key"`
	Value           string      `json:"value"`
	Target          []string    `json:"target"`
	ConfigurationID interface{} `json:"configurationId"`
	UpdatedAt       int64       `json:"updatedAt"`
	CreatedAt       int64       `json:"createdAt"`
}

type Handler struct {
	Api httpApi.API
}

func (h *Handler) Create(projectID string, env CreateOrUpdateEnv, teamId string) (string, error) {
	url := fmt.Sprintf("/v6/projects/%s/env", projectID)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request(http.MethodPost, url, env)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var createdEnv Env
	err = json.NewDecoder(res.Body).Decode(&createdEnv)
	if err != nil {
		return "", nil
	}
	log.Printf("%+v\n\n", createdEnv)

	return createdEnv.ID, nil
}

// Read returns environment variables associated with a project
func (h *Handler) Read(projectID string, teamId string) (envs []Env, err error) {
	url := fmt.Sprintf("/v6/projects/%s/env", projectID)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("GET", url, nil)
	if err != nil {
		return []Env{}, fmt.Errorf("Unable to fetch environment variables from vercel: %w", err)
	}
	defer res.Body.Close()

	// EnvResponse is only a subset of available data but all we care about
	// See https://vercel.com/docs/api#endpoints/projects/get-project-environment-variables
	type EnvResponse struct {
		Envs []Env `json:"envs"`
	}

	var envResponse EnvResponse

	err = json.NewDecoder(res.Body).Decode(&envResponse)
	if err != nil {
		return []Env{}, fmt.Errorf("Unable to unmarshal environment variables response: %w", err)
	}
	log.Printf("%+v\n\n", envResponse)
	return envResponse.Envs, nil
}
func (h *Handler) Update(projectID string, envID string, env CreateOrUpdateEnv, teamId string) error {
	url := fmt.Sprintf("/v6/projects/%s/env/%s", projectID, envID)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := h.Api.Request("PATCH", url, env)
	if err != nil {
		return fmt.Errorf("Unable to update env: %w", err)
	}
	defer res.Body.Close()
	return nil
}
func (h *Handler) Delete(projectID, envKey string, teamId string) error {
	url := fmt.Sprintf("/v4/projects/%s/env/%s", projectID, envKey)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete env: %w", err)
	}
	defer res.Body.Close()
	return nil
}
