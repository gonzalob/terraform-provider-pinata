provider "pinata" {
  root  = "http://localhost" # optional; defaults to pinata's live api
  token = "ey...d9"          # pinata jwk token with v3 files:write, legacy:pinning:pinFileToIPFS permissions
}
