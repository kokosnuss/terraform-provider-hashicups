terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
  required_version = ">= 1.1.0"
}

provider "hashicups" {
  username = "education"
  password = "test123"
  host     = "http://localhost:9090"
}

resource "hashicups_coffee" "edu" {
  name       = "Atruvia Terraform Boost"
  teaser     = "Double Espresso"
  collection = "New Arrivals"
  origin     = "Techdays 2025"
  price      = 150
  image      = "/terraform.png"
  ingredients = [{
    name     = "Espresso"
    quantity = 200
    unit     = "ml"
    },
    {
      name     = "Pumpkin Spice"
      quantity = 10
      unit     = "g"
    },
  ]
}

