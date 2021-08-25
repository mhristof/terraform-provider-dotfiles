package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
)

func config(file string) string {
	return fmt.Sprintf(heredoc.Doc(`
		resource "dotfiles_file" "dots" {
		  src = "%s"
		}
	`), file)

}
func TestFile(t *testing.T) {
	resourceName := "dotfiles_file.dots"
	file := "provider.go"
	srcAbs, err := filepath.Abs(file)
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	providers := map[string]*schema.Provider{
		"dotfiles": NewWithRoot(dir),
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		CheckDestroy: func(s *terraform.State) error {
			rs := s.RootModule().Resources[resourceName].Primary.ID
			if _, err := os.Stat(rs); !os.IsNotExist(err) {
				return errors.Errorf("file found [%s]", rs)
			}

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config(file),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						if _, err := os.Stat(rs.Primary.ID); os.IsNotExist(err) {
							return errors.Wrapf(err, "file [%s] not found", rs.Primary.ID)
						}

						return nil
					},
					resource.TestCheckResourceAttr(resourceName, "id", filepath.Join(dir, file)),
					resource.TestCheckResourceAttr(resourceName, "src_abs", srcAbs),
					resource.TestCheckResourceAttr(resourceName, "path_abs", filepath.Join(dir, file)),
					resource.TestCheckResourceAttr(resourceName, "src", file),
				),
			},
			{
				// idempotency
				Config:             config(file),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// file got removed from an outside source
				Config: config(file),
				PreConfig: func() {
					os.Remove(filepath.Join(dir, file))
				},
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						if _, err := os.Stat(rs.Primary.ID); os.IsNotExist(err) {
							return errors.Wrapf(err, "file [%s] not found", rs.Primary.ID)
						}

						return nil
					},
				),
			},
			// {
			// 	// file got replaced with another file
			// 	Config: config(file),
			// 	PreConfig: func() {
			// 		path := filepath.Join(dir, file)
			// 		os.Remove(path)
			// 		err := ioutil.WriteFile(path, []byte("test"), 0644)
			// 		if err != nil {
			// 			t.Fatal(err)
			// 		}

			// 	},
			// 	ExpectNonEmptyPlan: false,
			// },
		},
	})
}
