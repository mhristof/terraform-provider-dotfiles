package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinkCreate,
		Read:   resourceLinkRead,
		Update: resourceLinkUpdate,
		Delete: resourceLinkDelete,

		Schema: map[string]*schema.Schema{
			"dest": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"strip_source": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceLinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg)

	log.Println(fmt.Sprintf("[DEBUG] config: %+v", config))

	dest := d.Get("dest").(string)
	d.SetId(dest)
	source := d.Get("source").(string)

	_, err := os.Stat(source)

	if os.IsNotExist(err) {
		return err
	}

	abs, err := filepath.Abs(source)
	if err != nil {
		return err
	}

	lnDest := linkDest(config.root, dest, d.Get("strip_source").(bool))

	os.Symlink(abs, lnDest)

	log.Println(fmt.Sprintf("[DEBUG] ln %s %s", abs, lnDest))
	return resourceLinkRead(d, meta)
}

func linkDest(dir, name string, stripName bool) string {
	if stripName {
		name = path.Base(name)
	}
	return filepath.Join(dir, name)
}

func resourceLinkRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] reading resource event")
	return nil
}

func resourceLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO
	return nil
}

func resourceLinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg)
	source := d.Get("source").(string)

	dest := linkDest(config.root, source, d.Get("strip_source").(bool))

	os.Remove(dest)

	return nil
}
