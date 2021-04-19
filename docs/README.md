# Openstack Plugin

The Openstack Packer plugin provides a builder that is able to create new images
for use with OpenStack. The builder takes a source image, runs any provisioning
necessary on the image after launching it, then creates a new reusable image.
This reusable image can then be used as the foundation of new servers that are
launched within OpenStack. The builder will create temporary keypairs that
provide temporary access to the server while the image is being created. This
simplifies configuration quite a bit.

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    openstack = {
      version = ">= 0.0.1"
      source  = "github.com/hashicorp/openstack"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/hashicorp/packer-plugin-openstack/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


#### From Source

If you prefer to build the plugin from its source code, clone the GitHub
repository locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-openstack` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


## Plugin Contents

### Builder

- [builder](/docs/builders/openstack.mdx) - The Openstack Packer builder is able to create new images for use with OpenStack.