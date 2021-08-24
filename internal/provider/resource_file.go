package provider

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileCreate,
		ReadContext:   resourceFileRead,
		UpdateContext: resourceFileUpdate,
		DeleteContext: resourceFileDelete,
		Schema: map[string]*schema.Schema{
			"src": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func calculatePath(root, src, dest string) string {
	if dest == "" {
		dest = src
	}

	if strings.HasPrefix(dest, "/") {
		// dest is abs path
		return dest
	}

	return filepath.Join(root, dest)
}

func resourceFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*apiClient)

	path := calculatePath(config.root, d.Get("src").(string), d.Get("path").(string))

	abs, err := filepath.Abs(d.Get("src").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("ln", abs, path)

	os.Symlink(abs, path)

	d.SetId(path)
	return diags
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFileRead(ctx, d, m)
}

func resourceFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	os.Remove(d.Id())

	return diags
}
