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
  name   = "Double Coffee"
  teaser = "For my roMAN"
  price  = 150
  image  = "/terraform.png"
  ingredients = [{
    name     = "Espresso"
    quantity = 50
    unit     = "ml"
    },
  ]
}

