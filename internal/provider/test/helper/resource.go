package helper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestCheckNumberOfResources(expectedNum int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resources := s.RootModule().Resources
		l := len(resources)
		if l == expectedNum {
			return nil
		}
		return fmt.Errorf("invalid number of resources -> expected: %v, got: %v, current resources: %v", expectedNum, l, resources)
	}
}

func TestCheckResourceExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.RootModule().Resources[resource]
		if resourceState != nil {
			return nil
		}
		return fmt.Errorf("resource '%v' not found", resource)
	}
}

func TestCheckResourceNotExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.RootModule().Resources[resource]
		if resourceState == nil {
			return nil
		}
		return fmt.Errorf("resource '%v' found, but it should not exist. resource state: %v", resource, resourceState)
	}
}
