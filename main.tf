provider "dotfiles" {}



resource "dotfiles_link" "foo" {
  source = ".zshrc"
}
