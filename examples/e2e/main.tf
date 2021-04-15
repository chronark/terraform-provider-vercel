terraform {
  required_providers {
    vercel = {
      source  = "hashicorp.com/chronark/vercel"
      version = "9000.1"
    }
  }
}

provider "vercel" {
  token = "cwWSCeYIsBBlYuN9XdbHwKC8"
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


resource "vercel_env" "my_env" {
  project_id = vercel_project.my_project.id
  type       = "plain"
  key        = "hello"
  value      = "world"
  target     = ["production", "preview", "development"]
  team_id    = data.vercel_team.triebwork.id

}