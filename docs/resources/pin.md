---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pinata_pin Resource - pinata"
subcategory: ""
description: |-
  Manages an IPFS pin
---

# pinata_pin (Resource)

Manages an IPFS pin

## Example Usage

```terraform
resource "pinata_pin" "pin" {
  version = 0               # cid version to use. valid values are 0 and 1
  paths   = ["resource.tf"] # list of files (in the workspace) to pin to ipfs for this resource
  name    = "example"       # optional. the descriptive name on pinata for this pin
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `paths` (List of String) Local paths for the pin
- `version` (Number) The CID version to use

### Optional

- `name` (String) Resource name

### Read-Only

- `cid` (String) The pin's IPFS Content ID
- `hash` (String) Resource checksum
- `id` (String) Pinata ID for the pin
