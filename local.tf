terraform {
  required_providers {
    dotfiles = {
      source  = "github.com/mhristof/dotfiles-local"
      version = "0.1.0"
    }
  }
}

provider "dotfiles" {
  root = "/tmp/dotfiles"
  # example configuration here
}
