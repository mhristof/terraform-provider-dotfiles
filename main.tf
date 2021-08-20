provider "dotfiles" {
  root = "/tmp"
}

resource "dotfiles_curl" "foo" {
  url = "https://raw.githubusercontent.com/alphagov/collections/main/LICENCE.txt"
}
