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
		ResourcesMap: map[string]*schema.Resource{
			"dotfiles_link": resourceEvent(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// TODO
	return nil, nil
}

func resourceEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventCreate,
		Read:   resourceEventRead,
		Update: resourceEventUpdate,
		Delete: resourceEventDelete,

		Schema: map[string]*schema.Schema{
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type File struct {
	path string
}

func resourceEventCreate(d *schema.ResourceData, meta interface{}) error {
	source := d.Get("source").(string)
	d.SetId(source)
	return resourceEventRead(d, meta)
}

func resourceEventRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceEventUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO
	return nil
}

func resourceEventDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO
	return nil
}
