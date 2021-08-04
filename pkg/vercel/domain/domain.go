package domain

import (
	"encoding/json"
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
	"net/http"
)

type CreateDomain struct {
	Name string `json:"name"`
}
type Domain struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	ServiceType         string   `json:"serviceType"`
	NsVerifiedAt        int64    `json:"nsVerifiedAt"`
	TxtVerifiedAt       int64    `json:"txtVerifiedAt"`
	CdnEnabled          bool     `json:"cdnEnabled"`
	CreatedAt           int64    `json:"createdAt"`
	ExpiresAt           int64    `json:"expiresAt"`
	BoughtAt            int64    `json:"boughtAt"`
	TransferredAt       int64    `json:"transferredAt"`
	VerificationRecord  string   `json:"verificationRecord"`
	Verified            bool     `json:"verified"`
	Nameservers         []string `json:"nameservers"`
	IntendedNameservers []string `json:"intendedNameservers"`
	Creator             struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"creator"`
}

type Handler struct {
	Api httpApi.API
}

func (h *Handler) Create(name string, teamId string) (string, error) {
	url := "/v5/domains"
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request(http.MethodPost, url, CreateDomain{Name: name})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var createdDomain Domain
	err = json.NewDecoder(res.Body).Decode(&createdDomain)
	if err != nil {
		return "", nil
	}

	return createdDomain.ID, nil
}

// Read returns metadata about a domain
func (h *Handler) Read(domainName string, teamId string) (domain Domain, err error) {
	url := fmt.Sprintf("/v5/domains/%s", domainName)
	if teamId != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamId)
	}

	res, err := h.Api.Request("GET", url, nil)
	if err != nil {
		return Domain{}, fmt.Errorf("Unable to fetch domain from vercel: %w", err)
	}
	defer res.Body.Close()

	type GetDomainResponse struct {
		Domain Domain `json:"domain"`
	}
	getDomainResponse := GetDomainResponse{}
	err = json.NewDecoder(res.Body).Decode(&getDomainResponse)
	if err != nil {
		return Domain{}, fmt.Errorf("Unable to unmarshal domain response: %w", err)
	}

	return getDomainResponse.Domain, nil
}

func (h *Handler) Delete(domainName string, teamId string) error {
	url := fmt.Sprintf("/v5/domains/%s", domainName)
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
