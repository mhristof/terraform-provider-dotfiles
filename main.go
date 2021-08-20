package main

// https://www.hashicorp.com/blog/managing-google-calendar-with-terraform

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
			"dotfiles_link": resourceEvent(),
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

func resourceEventCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg)

	log.Println(fmt.Sprintf("[DEBUG] config: %+v", config))

	source := d.Get("source").(string)

	_, err := os.Stat(source)

	if os.IsNotExist(err) {
		return err
	}

	abs, err := filepath.Abs(source)
	if err != nil {
		return err
	}

	dest := filepath.Join(config.root, source)

	os.Symlink(abs, dest)

	log.Println(fmt.Sprintf("[DEBUG] ln %s %s", abs, dest))
	d.SetId(source)
	return resourceEventRead(d, meta)
}

func resourceEventRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] reading resource event")
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
