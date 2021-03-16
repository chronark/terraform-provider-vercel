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



data "vercel_user" "user" {}