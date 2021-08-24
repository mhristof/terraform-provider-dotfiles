---
page_title: "dotfiles_data_source Data Source - terraform-provider-dotfiles"
subcategory: ""
description: |-
  Sample data source in the Terraform provider dotfiles.
---

# Data Source `dotfiles_data_source`

Sample data source in the Terraform provider dotfiles.

## Example Usage

```terraform
data "dotfiles_data_source" "example" {
  sample_attribute = "foo"
}
```

## Schema

### Required

- **sample_attribute** (String, Required) Sample attribute.

### Optional

- **id** (String, Optional) The ID of this resource.


