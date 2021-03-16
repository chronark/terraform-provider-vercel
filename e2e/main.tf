terraform {
  required_providers {
    fauna = {
      source  = "hashicorp.com/chronark/vercel"
      version = "9000.1"
    }
  }
}

provider "fauna" {
  fauna_key = "fnAEEWZgLrACB6LzAzbAotOEPVrqCQKX1-rbedfw"
}




resource "fauna_collection" "my_collection" {
  name = "terraform"
}

resource "fauna_index" "my_index" {
  sources = [fauna_collection.my_collection.name]
  name    = "my_index"
}

