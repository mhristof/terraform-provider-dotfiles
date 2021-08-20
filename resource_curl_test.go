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

func calculateSHA256(t *testing.T, file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func TestCurl(t *testing.T) {
	var cases = []struct {
		name     string
		fs       map[string]string
		validate func(*testing.T, string, string, string)
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
			validate: func(t *testing.T, dir, file, name string) {
				assert.Equal(
					t,
					"4bf67172f2ada15c5538c37f86ed157300171540b92e119ab59cf6b1a2cb48b7",
					calculateSHA256(t, filepath.Join(dir, file)),
					name,
				)
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
			validate: func(t *testing.T, dir, file, name string) {
				assert.Equal(
					t,
					"4bf67172f2ada15c5538c37f86ed157300171540b92e119ab59cf6b1a2cb48b7",
					calculateSHA256(t, filepath.Join(dir, file)),
					name,
				)
			},
		},
		{
			name: "zip file with binary extraction",
			fs: map[string]string{
				"main.tf": heredoc.Doc(`
					provider "dotfiles" {}

					resource "dotfiles_curl" "foo" {
  						url     = "https://github.com/terraform-linters/tflint/releases/download/v0.31.0/tflint_darwin_amd64.zip"
  						extract = "tflint"
					}

					output "file" {
						value = "tflint"
					}
				`),
			},
			validate: func(t *testing.T, dir, file, name string) {
				assert.Equal(
					t,
					"4bf67172f2ada15c5538c37f86ed157300171540b92e119ab59cf6b1a2cb48b7",
					calculateSHA256(t, filepath.Join(dir, file)),
					name,
				)
			},
		},
		{
			name: "tar.gz file with binary extraction",
			fs: map[string]string{
				"main.tf": heredoc.Doc(`
					provider "dotfiles" {}

					resource "dotfiles_curl" "foo" {
						url     = "https://github.com/golangci/golangci-lint/releases/download/v1.42.0/golangci-lint-1.42.0-darwin-amd64.tar.gz"
						extract = "golangci-lint"
					}

					output "file" {
						value = "golangci-lint"
					}
				`),
			},
			validate: func(t *testing.T, dir, file, name string) {
				assert.Equal(
					t,
					"e34fd545ed520ad1e62abb894ca5b858b23b06110dfd19f3ca11397414c79385",
					calculateSHA256(t, filepath.Join(dir, file)),
					name,
				)
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
			test.validate(t, dir, terraform.Output(t, terraformOptions, "file"), test.name)
		})
	}
}
