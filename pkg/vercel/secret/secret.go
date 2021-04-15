package secret

import (
	"encoding/json"
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
	"time"
)

type Secret struct {
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	TeamID    string    `json:"teamId"`
	UserID    string    `json:"userId"`
	Created   time.Time `json:"created"`
	CreatedAt int64     `json:"createdAt"`
}

type CreateSecret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Handler struct {
	Api httpApi.API
}

func (h *Handler) Create(secret CreateSecret) (string, error) {
	res, err := h.Api.Request("POST", "/v2/now/secrets", secret)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var createdSecret Secret
	err = json.NewDecoder(res.Body).Decode(&createdSecret)
	if err != nil {
		return "", nil
	}

	return createdSecret.UID, nil
}

// Read returns environment variables associated with a project
func (h *Handler) Read(secretID, teamId string) (secret Secret, err error) {
	url := fmt.Sprintf("/v3/now/secrets/%s", secretID)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := h.Api.Request("GET", url, nil)
	if err != nil {
		return Secret{}, fmt.Errorf("Unable to fetch secret from vercel: %w", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&secret)
	if err != nil {
		return Secret{}, fmt.Errorf("Unable to unmarshal environment variables response: %w", err)
	}
	return secret, nil
}
func (h *Handler) Update(oldName, newName, teamId string) error {
	url := fmt.Sprintf("/v2/now/secrets/%s", oldName)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	payload := struct {
		Name string `json:"name"`
	}{
		Name: newName,
	}

	res, err := h.Api.Request("PATCH", url, payload)
	if err != nil {
		return fmt.Errorf("Unable to update secret: %w", err)
	}
	defer res.Body.Close()
	return nil
}
func (h *Handler) Delete(secretName, teamId string) error {
	url := fmt.Sprintf("/v2/now/secrets/%s", secretName)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete secret: %w", err)
	}
	defer res.Body.Close()
	return nil
}
