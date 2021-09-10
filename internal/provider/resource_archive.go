package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mhristof/go-archive"
)

func resourceArchive() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchiveCreate,
		ReadContext:   resourceArchiveRead,
		UpdateContext: resourceArchiveUpdate,
		DeleteContext: resourceArchiveDelete,
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dest": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sha256sum": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArchiveCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*apiClient)

	file := d.Get("file").(string)

	zip := archive.NewURL(d.Get("url").(string))

	data, err := zip.ExtractFile(file)

	log.Println("extracted", file)

	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "File not found",
			Detail:   fmt.Sprintf("File [%s] not found in the archive", file),
		})
	}

	dest := d.Get("dest").(string)
	if dest == "" {
		dest = filepath.Join(config.root, path.Base(file))
	}

	log.Println(fmt.Sprintf("creating file dest: %+v", dest))

	err = os.WriteFile(dest, data, 0755)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot create file",
			Detail:   fmt.Sprintf("Unable to create file [%s], %+v", dest, err),
		})
	}

	d.Set("sha256sum", fmt.Sprintf("%x", sha256.Sum256(data)))
	d.SetId(dest)
	return diags
}

func resourceArchiveRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	data, err := ioutil.ReadFile(d.Id())
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%+v", err),
			Detail:   fmt.Sprintf("File [%s] not found in the archive", d.Id()),
		})
	}

	if fmt.Sprintf("%x", sha256.Sum256(data)) != d.Get("sha256sum").(string) {
		d.SetId("")
	}

	return diags
}

func resourceArchiveUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceArchiveRead(ctx, d, m)
}

func resourceArchiveDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	os.Remove(d.Id())
	return diags
}
