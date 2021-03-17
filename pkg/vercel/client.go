package vercel

import (
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/user"
)

type Client struct {
	Project *project.ProjectHandler
	User    *user.UserHandler
}

func New(token string) *Client {
	api := httpApi.New(token)

	return &Client{
		Project: &project.ProjectHandler{
			Api: api,
		},
		User: &user.UserHandler{
			Api: api,
		},
	}
}
