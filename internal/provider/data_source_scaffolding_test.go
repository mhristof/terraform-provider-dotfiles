package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDotfiles(t *testing.T) {
	t.Skip("data source not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDotfiles,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.dotfiles_data_source.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccDataSourceDotfiles = `
data "dotfiles_data_source" "foo" {
  sample_attribute = "bar"
}
`
