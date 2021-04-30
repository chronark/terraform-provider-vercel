package project

import (
	"encoding/json"
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type ProjectHandler struct {
	Api httpApi.API
}

func (p *ProjectHandler) Create(project CreateOrUpdateProject, teamId string) (string, error) {
	url := "/v6/projects"
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("POST", url, project)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var createdProject Project
	err = json.NewDecoder(res.Body).Decode(&createdProject)
	if err != nil {
		return "", nil
	}

	return createdProject.ID, nil
}
func (p *ProjectHandler) Read(id string, teamId string) (project Project, err error) {
	url := fmt.Sprintf("/v1/projects/%s", id)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("GET", url, nil)
	if err != nil {
		return Project{}, fmt.Errorf("Unable to fetch project from vercel: %w", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return Project{}, fmt.Errorf("Unable to unmarshal project: %w", err)
	}
	return project, nil
}
func (p *ProjectHandler) Update(id string, project CreateOrUpdateProject, teamId string) error {
	url := fmt.Sprintf("/v2/projects/%s", id)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("PATCH", url, project)
	if err != nil {
		return fmt.Errorf("Unable to update project: %w", err)
	}
	defer res.Body.Close()
	return nil
}
func (p *ProjectHandler) Delete(id string, teamId string) error {
	url := fmt.Sprintf("/v1/projects/%s", id)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete project: %w", err)
	}
	defer res.Body.Close()
	return nil
}
