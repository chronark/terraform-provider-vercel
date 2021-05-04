package vercel

import (
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/env"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/alias"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/dns"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/domain"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/secret"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/team"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/user"
)

type Client struct {
	Project *project.ProjectHandler
	User    *user.UserHandler
	Env     *env.Handler
	Secret  *secret.Handler
	Team    *team.Handler
	Alias   *alias.Handler
	Domain  *domain.Handler
	DNS     *dns.Handler
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
		Env:    &env.Handler{Api: api},
		Secret: &secret.Handler{Api: api},
		Team:   &team.Handler{Api: api},
		Alias:  &alias.Handler{Api: api},
		Domain: &domain.Handler{Api: api},
		DNS:    &dns.Handler{Api: api},
	}
}
