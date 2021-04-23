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



resource "vercel_project" "my_project" {
  name = "test"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }


}


resource "vercel_domain" "chronark-com" {
  name = "chronark.com"
}