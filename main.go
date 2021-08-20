package main

// https://www.hashicorp.com/blog/managing-google-calendar-with-terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"root": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "./",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dotfiles_link": resourceLink(),
		},
		ConfigureFunc: providerConfigure,
	}
}

type cfg struct {
	root string
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return &cfg{
		root: d.Get("root").(string),
	}, nil
}
