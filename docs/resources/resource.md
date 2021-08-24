---
page_title: "dotfiles_resource Resource - terraform-provider-dotfiles"
subcategory: ""
description: |-
  Sample resource in the Terraform provider dotfiles.
---

# Resource `dotfiles_resource`

Sample resource in the Terraform provider dotfiles.

## Example Usage

```terraform
resource "dotfiles_resource" "example" {
  sample_attribute = "foo"
}
```

## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **sample_attribute** (String, Optional) Sample attribute.


