package provider

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func config() string {
	return heredoc.Doc(`
		resource "dotfiles_file" "dots" {
		  src = "local.tf"
		}
	`)

}
func TestAccAWSS3Bucket_Tags_withNoSystemTags(t *testing.T) {
	resourceName := "dotfiles_file.dots"

	resource.ParallelTest(t, resource.TestCase{
		//PreCheck:     func() { testAccPreCheck(t) },
		//ErrorCheck:   testAccErrorCheck(t, s3.EndpointsID),
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config(),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckAWSS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "asdf"),
				//resource.TestCheckResourceAttr(resourceName, "tags.Key1", "AAA"),
				//resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
				//resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
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
