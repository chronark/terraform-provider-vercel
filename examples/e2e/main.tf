terraform {
  required_providers {
    vercel = {
      source  = "hashicorp.com/chronark/vercel"
      version = "9000.1"
    }
  }
}

provider "vercel" {
}


data "vercel_team" "triebwork" {
  slug = "triebwork"
}

resource "vercel_project" "my_project" {
  name = "mercury-via-terraform"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }
  team_id = data.vercel_team.triebwork.id
}


resource "vercel_secret" "my_secret" {
  name    = "hello"
  value   = "world"
  team_id = data.vercel_team.triebwork.id
}

resource "vercel_env" "my_env" {
  team_id    = data.vercel_team.triebwork.id
  project_id = vercel_project.my_project.id

  type   = "secret"
  value  = vercel_secret.my_secret.id
  key    = "world"
  target = ["production", "preview", "development"]

}