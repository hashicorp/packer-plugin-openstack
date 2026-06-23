# Copyright IBM Corp. 2013, 2026
# SPDX-License-Identifier: MPL-2.0

# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "OpenStack"
  description = "The OpenStack multi-component plugin can be used with HashiCorp Packer to create custom images."
  identifier = "packer/hashicorp/openstack"
  flags = ["hcp-ready"]
  component {
    type = "builder"
    name = "OpenStack"
    slug = "openstack"
  }
}
