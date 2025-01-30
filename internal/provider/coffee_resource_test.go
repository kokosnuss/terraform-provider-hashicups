package provider

import (
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
			// ImportState testing
			// Update and Read testing
			// Delete testing automatically occurs in TestCase
		},
	})
}
