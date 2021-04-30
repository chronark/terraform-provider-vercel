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
  name = "hallo-jannik"
  # git_repository {
  #   type = "github"
  #   repo = "chronark/flare"
  # }
  team_id = data.vercel_team.triebwork.id
}

