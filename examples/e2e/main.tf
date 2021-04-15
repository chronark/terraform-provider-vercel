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

data "vercel_team" "trieb_work" {
  slug = "triebwork"
}


resource "vercel_project" "dev-website" {
  name = "dev-website"
  git_repository {
    type = "github"
    repo = "trieb.work/dev-website"
  }
  framework = "nextjs"
  // team_id   = data.vercel_team.trieb_work.id
}