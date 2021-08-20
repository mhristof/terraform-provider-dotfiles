package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCurl() *schema.Resource {
	return &schema.Resource{
		Create: resourceCurlCreate,
		Read:   resourceCurlRead,
		Update: resourceCurlUpdate,
		Delete: resourceCurlDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dest": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCurlCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg)

	log.Println(fmt.Sprintf("[DEBUG] config: %+v", config))

	url, dest := getDest(d, meta)

	d.SetId(dest)

	Wget(url, dest)

	log.Println(fmt.Sprintf("[DEBUG] wget %s %s", url, dest))
	return resourceCurlRead(d, meta)
}

func Wget(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getDest(d *schema.ResourceData, meta interface{}) (string, string) {
	config := meta.(*cfg)

	url := d.Get("url").(string)
	dest := d.Get("dest").(string)
	if dest == "" {
		dest = path.Base(url)
	}

	return url, filepath.Join(config.root, dest)
}

func resourceCurlRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] reading resource event")
	return nil
}

func resourceCurlUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO
	return nil
}

func resourceCurlDelete(d *schema.ResourceData, meta interface{}) error {
	_, dest := getDest(d, meta)
	os.Remove(dest)
	return nil
}
