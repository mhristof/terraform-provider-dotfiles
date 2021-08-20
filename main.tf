provider "dotfiles" {
  root = "/tmp"
}

resource "dotfiles_link" "foo" {
  source = ".zshrc"
}
