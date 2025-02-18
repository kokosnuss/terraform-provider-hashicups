package provider

import (
	"terraform-provider-hashicups/internal/provider/test/helper"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCoffeeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "hashicups_coffee" "test" {
				name = "terraspiced latte"
				teaser = "exclusively for techdays 2025"
				price = 150
				image = "/terraform.png"
				ingredients = [{
					name = "Espresso"
					quantity = 50
					unit = "ml"
					},
					{
					name = "Steamed Milk"
					quantity = 100
					unit = "ml"
				}]

				}
				resource "hashicups_coffee" "second_test" {
				name = "random_mix"
				teaser = "test only, not for consumption"
				price = -1
				image = "/terraform.png"
				ingredients = [{
					name = "Hot Water"
					quantity = 1
					unit = "l"
					}]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					helper.TestCheckNumberOfResources(2),
					helper.TestCheckResourceExists("hashicups_coffee.test"),
					helper.TestCheckResourceExists("hashicups_coffee.second_test"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "name", "terraspiced latte"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.0.name", "Espresso"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.0.quantity", "50"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.0.unit", "ml"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.1.name", "Steamed Milk"),
				),
			},
			// ImportState testing
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "hashicups_coffee" "test" {
				name = "terraspiced coffein booster"
				teaser = "exclusively for techdays 2025"
				price = 250
				image = "/terraform.png"
				ingredients = [{
					name = "Espresso"
					quantity = 10
					unit = "dl"
					}]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					helper.TestCheckNumberOfResources(1),
					helper.TestCheckResourceExists("hashicups_coffee.test"),
					helper.TestCheckResourceNotExists("hashicups_coffee.second_test"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "name", "terraspiced coffein booster"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.0.name", "Espresso"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.0.quantity", "10"),
					resource.TestCheckResourceAttr("hashicups_coffee.test", "ingredients.0.unit", "dl"),
					resource.TestCheckNoResourceAttr("hashicups_coffee.test", "ingredients.1.name"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
