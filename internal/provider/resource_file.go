package provider

import (
	"context"
	"fmt"
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
			"src_abs": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"path_abs": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "File already exists",
			Detail:   fmt.Sprintf("File [%s] exists in the filesystem and terraform would overwrite it. Please remove the file and try again", path),
		})
	}
	os.Symlink(abs, path)

	d.Set("src_abs", abs)
	d.Set("path_abs", path)
	d.SetId(path)
	return diags
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	if _, err := os.Stat(id); os.IsNotExist(err) {
		d.SetId("")
		return diags
	}

	fileInfo, err := os.Lstat(id)
	if err != nil {
		// TODO: test
		return diag.FromErr(err)
	}

	if fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
		// TODO: test
		d.SetId("")
		return diags
	}

	dest, err := os.Readlink(id)
	if err != nil {
		// TODO: test
		return diag.FromErr(err)
	}

	if dest != d.Get("src_abs").(string) {
		// TODO: test
		d.SetId("")
		return diags
	}

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
