package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func createFs(t *testing.T, files map[string]string) (string, func()) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		t.Fatal(err)
	}

	for name, content := range files {
		err = ioutil.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	return dir, func() {
		os.RemoveAll(dir)
	}
}

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

		terraformD, err := filepath.Abs("./terraform.d")
		if err != nil {
			t.Fatal(err)
		}

		providerTf, err := filepath.Abs("./providers.tf")
		if err != nil {
			t.Fatal(err)
		}

		os.Symlink(terraformD, fmt.Sprintf("%s/terraform.d", dir))
		os.Symlink(providerTf, fmt.Sprintf("%s/providers.tf", dir))

		fmt.Println(fmt.Sprintf("dir: %+v", dir))

		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: dir,
		})

		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		assert.Nil(t, test.validate(dir), test.name)
	}
}
