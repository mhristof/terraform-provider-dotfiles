provider "dotfiles" {
  root = "/tmp"
}

resource "dotfiles_curl" "foo" {
  url     = "https://github.com/golangci/golangci-lint/releases/download/v1.42.0/golangci-lint-1.42.0-darwin-amd64.tar.gz"
  extract = "golangci-lint"
  url     = "https://github.com/terraform-linters/tflint/releases/download/v0.31.0/tflint_darwin_amd64.zip"
  extract = "tflint"
}
