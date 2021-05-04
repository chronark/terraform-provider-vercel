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


resource "vercel_domain" "chronark_com" {
  name = "chronark.com"
}


resource "vercel_dns" "www" {
  domain = vercel_domain.chronark_com.name
  type   = "CNAME"
  value  = "www.${vercel_domain.chronark_com.name}"
  name   = "www"
}