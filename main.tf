provider "dotfiles" {
  root = "/tmp"
}

resource "dotfiles_link" "foo" {
  dest         = ".zshrc"
  source       = "fixtures/.zshrc"
  strip_source = true
}
