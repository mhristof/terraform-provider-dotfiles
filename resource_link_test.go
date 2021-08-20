package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestLink(t *testing.T) {
	var cases = []struct {
		name     string
		fs       map[string]string
		validate func(string) error
	}{
		{
			name: "simple link creation",
			fs: map[string]string{
				"zshrc": heredoc.Doc(`
					this file is intentionally left blank
				`),
				"main.tf": heredoc.Doc(`
					provider "dotfiles" {}

					resource "dotfiles_link" "foo" {
					  dest         = ".zshrc"
					  source       = "zshrc"
					}
				`),
			},
			validate: func(dir string) error {
				if _, err := os.Stat(filepath.Join(dir, ".zshrc")); os.IsNotExist(err) {
					return err
				}

				return nil
			},
		},
	}

	for _, test := range cases {
		dir, cleanup := createFs(t, test.fs)
		defer cleanup()

		fmt.Println(fmt.Sprintf("dir: %+v", dir))

		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: dir,
		})

		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		assert.Nil(t, test.validate(dir), test.name)
	}
}
