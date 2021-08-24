terraform {
  required_providers {
    dotfiles = {
      source  = "github.com/mhristof/dotfiles-local"
      version = "0.1.0"
    }
  }
}

provider "dotfiles" {
  # example configuration here
}
