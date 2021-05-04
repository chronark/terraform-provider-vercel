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



resource "vercel_project" "flare_tf" {
  name = "flare-tf"
  git_repository {
    type = "github"
    repo = "chronark/flare"
  }


}


resource "vercel_domain" "flare_chronark_com" {
  name = "chronark.com"
}


resource "vercel_dns" "flare" {
  domain = vercel_domain.flare_chronark_com.name
  type   = "CNAME"
  value  = "flare.${vercel_domain.flare_chronark_com.name}"
  name   = "flare"
}


resource "vercel_alias" "flare" {
  project_id = vercel_project.flare_tf.id
  domain     = "flare.chronark.com"
}