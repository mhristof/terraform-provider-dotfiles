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

resource "dotfiles_file" "dots" {
  src = "README.md"
}

output "dots" {
  value = dotfiles_file.dots
}
