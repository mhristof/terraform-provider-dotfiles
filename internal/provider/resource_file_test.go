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
	dir, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	providers := map[string]*schema.Provider{
		"dotfiles": NewWithRoot(dir),
	}

	resource.ParallelTest(t, resource.TestCase{
		//PreCheck:     func() { testAccPreCheck(t) },
		//ErrorCheck:   testAccErrorCheck(t, s3.EndpointsID),
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
				),
			},
			// {
			// 	ResourceName:            resourceName,
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			// },
			// {
			// 	Config: testAccAWSS3BucketConfig_withUpdatedTags(bucketName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAWSS3BucketExists(resourceName),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key3", "XXX"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key5", "EEE"),
			// 	),
			// },
			// {
			// 	Config: testAccAWSS3BucketConfig_withNoTags(bucketName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAWSS3BucketExists(resourceName),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
			// 	),
			// },
			// Verify update from 0 tags.
			// {
			// 	Config: testAccAWSS3BucketConfig_withTags(bucketName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAWSS3BucketExists(resourceName),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key1", "AAA"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
			// 	),
			// },
		},
	})
}
