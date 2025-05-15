resource "pinata_pin" "pin" {
  version = 0               # cid version to use. valid values are 0 and 1
  paths   = ["resource.tf"] # list of files (in the workspace) to pin to ipfs for this resource
  name    = "example"       # optional. the descriptive name on pinata for this pin
}
