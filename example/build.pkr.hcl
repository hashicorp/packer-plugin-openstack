
source "openstack" "example" {
  domain_name       = "Default"
  flavor            = "m1.tiny"
  identity_endpoint = "http://<devstack-ip>:5000/v3"
  image_name        = "Test image"
  insecure          = "true"
  password          = "<your admin password>"
  region            = "RegionOne"
  source_image      = "<image id>"
  ssh_username      = "root"
  tenant_name       = "admin"
  username          = "admin"
}

build {
  sources = ["source.openstack.example"]
}
