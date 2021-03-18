terraform {
  required_providers {
    vercel = {
      source  = "hashicorp.com/chronark/vercel"
      version = "9000.1"
    }
  }
}

provider "vercel" {
  token = "wsByP9ptGqn7snGvvY00aDzn"
}

data "vercel_user" "me" {}

resource "vercel_project" "my_project" {
  name = "mercury-via-terraform"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }



}

resource "vercel_env" "env" {
  project_id = vercel_project.my_project.id
  type       = "plain"
  key        = "key"
  value      = "value"
  target     = ["production", "preview", "development"]
}