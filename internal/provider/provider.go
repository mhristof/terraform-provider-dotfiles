package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func NewWithRoot(root string) *schema.Provider {
	provider := New()

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return &apiClient{
			root: root,
		}, nil
	}
	return provider
}

func New() *schema.Provider {
	pwd, _ := os.Getwd()
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"root": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  pwd,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "dotfiles_data_source": dataSourceDotfiles(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dotfiles_file":    resourceFile(),
			"dotfiles_archive": resourceArchive(),
		},
		ConfigureFunc: providerConfigure,
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	root := d.Get("root").(string)
	return &apiClient{
		root: root,
	}, nil
}

type apiClient struct {
	root string
}
