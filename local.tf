terraform {
  required_providers {
    dotfiles = {
      source  = "github.com/mhristof/dotfiles-local"
      version = "0.1.0"
    }
  }
}

provider "dotfiles" {
  root = "/tmp"
}

resource "dotfiles_archive" "archives" {
  url  = "https://github.com/terraform-linters/tflint/releases/download/v0.31.0/tflint_darwin_amd64.zip"
  file = "tflint"
}

resource "dotfiles_file" "dots" {
  src = "README.md"
}

output "archives" {
  value = dotfiles_archive.archives
}
