package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCurl(t *testing.T) {
	var cases = []struct {
		name     string
		fs       map[string]string
		validate func(string) error
	}{
		{
			name: "simple url with binary name",
			fs: map[string]string{
				"main.tf": heredoc.Doc(`
					provider "dotfiles" {}

					locals {
						dest = "foo.txt"
					}

					resource "dotfiles_curl" "foo" {
					  url = "https://raw.githubusercontent.com/alphagov/collections/main/LICENCE.txt"
					  dest = local.dest
					}

					output "file" {
						value = local.dest
					}
				`),
			},
		},
		{
			name: "simple url without binary name",
			fs: map[string]string{
				"main.tf": heredoc.Doc(`
					provider "dotfiles" {}

					resource "dotfiles_curl" "foo" {
					  url = "https://raw.githubusercontent.com/alphagov/collections/main/LICENCE.txt"
					}

					output "file" {
						value = "LICENCE.txt"
					}
				`),
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dir, cleanup := createFs(t, test.fs)
			defer cleanup()

			fmt.Println(fmt.Sprintf("dir: %+v", dir))

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: dir,
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApply(t, terraformOptions)

			file := terraform.Output(t, terraformOptions, "file")
			data, err := ioutil.ReadFile(filepath.Join(dir, file))
			if err != nil {
				t.Fatal(err)
			}

			hash := sha256.Sum256(data)
			assert.Equal(
				t,
				"4bf67172f2ada15c5538c37f86ed157300171540b92e119ab59cf6b1a2cb48b7",
				hex.EncodeToString(hash[:]),
				test.name,
			)

		})
	}
}
