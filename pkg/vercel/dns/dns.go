package dns

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type CreateRecord struct {
	// The type of record, it could be any valid DNS record.
	// See https://vercel.com/docs/api#endpoints/dns for details
	Type string `json:"type"`

	// A subdomain name or an empty string for the root domain.
	Name string `json:"name"`

	// The record value.
	Value string `json:"value"`

	// The TTL value. Must be a number between 60 and 2147483647. Default value is 60.
	TTL int `json:"ttl"`
}

type Record struct {
	// The unique ID of the DNS record. Always prepended with rec_.
	Id string `json:"id"`

	// The type of record, it could be any valid DNS record.
	// See https://vercel.com/docs/api#endpoints/dns for details
	Type string `json:"type"`

	// A subdomain name or an empty string for the root domain.
	Name string `json:"name"`

	// The record value.
	Value string `json:"value"`

	// The ID of the user who created the record or system if the record is an automatic record.
	Creator string `json:"creator"`

	// The date when the record was created.
	Created int `json:"created"`

	// The date when the record was updated.
	Updated int `json:"updated"`

	// The date when the record was created in milliseconds since the UNIX epoch.
	CreatedAt int `json:"createdAt"`

	// The date when the record was updated in milliseconds since the UNIX epoch.
	UpdatedAt int `json:"updatedAt"`
}

type Handler struct {
	Api httpApi.API
}

func (h *Handler) Create(domain string, record CreateRecord, teamId string) (string, error) {
	url := fmt.Sprintf("/v2/domains/%s/records", domain)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request(http.MethodPost, url, record)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	type CreateResponse struct {
		UID string `json:"uid"`
	}

	var createResponse CreateResponse

	err = json.NewDecoder(res.Body).Decode(&createResponse)
	if err != nil {
		return "", err
	}

	return createResponse.UID, nil
}

func (h *Handler) Read(domain, recordId, teamId string) (Record, error) {
	url := fmt.Sprintf("/v2/domains/%s/records?limit=1000", domain)
	if teamId != "" {
		url = fmt.Sprintf("%s&teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("GET", url, nil)
	if err != nil {
		return Record{}, fmt.Errorf("Unable to fetch dns records: %w", err)
	}
	defer res.Body.Close()

	type Response struct {
		Records []Record `json:"records"`
	}
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return Record{}, err
	}

	for _, record := range response.Records {
		if record.Id == recordId {
			return record, nil
		}
	}
	return Record{}, fmt.Errorf("Record with id %s was not found", recordId)

}

func (h *Handler) Delete(domain, recordId string, teamId string) error {
	url := fmt.Sprintf("/v2/domains/%s/records/%s", domain, recordId)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := h.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete dns record: %w", err)
	}
	defer res.Body.Close()
	return nil
}
