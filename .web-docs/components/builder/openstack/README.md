Type: `openstack`
Artifact BuilderId: `mitchellh.openstack`

The `openstack` Packer builder is able to create new images for use with
[OpenStack](http://www.openstack.org). The builder takes a source image, runs
any provisioning necessary on the image after launching it, then creates a new
reusable image. This reusable image can then be used as the foundation of new
servers that are launched within OpenStack. The builder will create temporary
keypairs that provide temporary access to the server while the image is being
created. This simplifies configuration quite a bit.

The builder does _not_ manage images. Once it creates an image, it is up to you
to use it or delete it.

~> **Note:** To use OpenStack builder with the OpenStack Newton (Oct 2016)
or earlier, we recommend you use Packer v1.1.2 or earlier version.

~> **OpenStack Liberty or later requires OpenSSL!** To use the OpenStack
builder with OpenStack Liberty (Oct 2015) or later you need to have OpenSSL
installed _if you are using temporary key pairs_, i.e. don't use
[`ssh_keypair_name`](#ssh_keypair_name) nor
[`ssh_password`](#ssh_password). All major
OS'es have OpenSSL installed by default except Windows. This have been resolved
in OpenStack Ocata(Feb 2017).

~> **Note:** OpenStack Block Storage volume support is available only for
V3 Block Storage API. It's available in OpenStack since Mitaka release (Apr
2016).

## Configuration Reference

There are many configuration options available for the builder. They are
segmented below into two categories: required and optional parameters. Within
each category, the available configuration keys are alphabetized.

In addition to the options listed here, a
[communicator](/packer/docs/templates/legacy_json_templates/communicator) can be configured for this
builder.

### Required:

<!-- Code generated from the comments of the AccessConfig struct in builder/openstack/access_config.go; DO NOT EDIT MANUALLY -->

- `username` (string) - The username or id used to connect to the OpenStack service. If not
  specified, Packer will use the environment variable OS_USERNAME or
  OS_USERID, if set. This is not required if using access token or
  application credential instead of password, or if using cloud.yaml.

- `password` (string) - The password used to connect to the OpenStack service. If not specified,
  Packer will use the environment variables OS_PASSWORD, if set. This is
  not required if using access token or application credential instead of
  password, or if using cloud.yaml.

- `identity_endpoint` (string) - The URL to the OpenStack Identity service. If not specified, Packer will
  use the environment variables OS_AUTH_URL, if set. This is not required
  if using cloud.yaml.

<!-- End of code generated from the comments of the AccessConfig struct in builder/openstack/access_config.go; -->


<!-- Code generated from the comments of the ImageConfig struct in builder/openstack/image_config.go; DO NOT EDIT MANUALLY -->

- `image_name` (string) - The name of the resulting image.

<!-- End of code generated from the comments of the ImageConfig struct in builder/openstack/image_config.go; -->


<!-- Code generated from the comments of the RunConfig struct in builder/openstack/run_config.go; DO NOT EDIT MANUALLY -->

- `source_image` (string) - The ID or full URL to the base image to use. This is the image that will
  be used to launch a new server and provision it. Unless you specify
  completely custom SSH settings, the source image must have cloud-init
  installed so that the keypair gets assigned properly.

- `source_image_name` (string) - The name of the base image to use. This is an alternative way of
  providing source_image and only either of them can be specified.

- `external_source_image_url` (string) - The URL of an external base image to use. This is an alternative way of
  providing source_image and only either of them can be specified.

- `source_image_filter` (ImageFilter) - Filters used to populate filter options. Example:
  
  ```json
  {
      "source_image_filter": {
          "filters": {
              "name": "ubuntu-16.04",
              "visibility": "protected",
              "owner": "d1a588cf4b0743344508dc145649372d1",
              "tags": ["prod", "ready"],
              "properties": {
                  "os_distro": "ubuntu"
              }
          },
          "most_recent": true
      }
  }
  ```
  
  This selects the most recent production Ubuntu 16.04 shared to you by
  the given owner. NOTE: This will fail unless *exactly* one image is
  returned, or `most_recent` is set to true. In the example of multiple
  returned images, `most_recent` will cause this to succeed by selecting
  the newest image of the returned images.
  
  -   `filters` (map of strings) - filters used to select a
  `source_image`.
      NOTE: This will fail unless *exactly* one image is returned, or
      `most_recent` is set to true. Of the filters described in
      [ImageService](https://developer.openstack.org/api-ref/image/v2/), the
      following are valid:
  
      -   name (string)
      -   owner (string)
      -   tags (array of strings)
      -   visibility (string)
      -   properties (map of strings to strings) (fields that can be set
          with `openstack image set --property key=value`)
  
  -   `most_recent` (boolean) - Selects the newest created image when
  true.
      This is most useful for selecting a daily distro build.
  
  You may set use this in place of `source_image` If `source_image_filter`
  is provided alongside `source_image`, the `source_image` will override
  the filter. The filter will not be used in this case.

- `flavor` (string) - The ID, name, or full URL for the desired flavor for the server to be
  created.

<!-- End of code generated from the comments of the RunConfig struct in builder/openstack/run_config.go; -->


### Optional:

<!-- Code generated from the comments of the AccessConfig struct in builder/openstack/access_config.go; DO NOT EDIT MANUALLY -->

- `user_id` (string) - Sets username

- `tenant_id` (string) - The tenant ID or name to boot the instance into. Some OpenStack
  installations require this. If not specified, Packer will use the
  environment variable OS_TENANT_NAME or OS_TENANT_ID, if set. Tenant is
  also called Project in later versions of OpenStack.

- `tenant_name` (string) - Tenant Name

- `domain_id` (string) - Domain ID

- `domain_name` (string) - The Domain name or ID you are authenticating with. OpenStack
  installations require this if identity v3 is used. Packer will use the
  environment variable OS_DOMAIN_NAME or OS_DOMAIN_ID, if set.

- `insecure` (bool) - Whether or not the connection to OpenStack can be done over an insecure
  connection. By default this is false.

- `region` (string) - The name of the region, such as "DFW", in which to launch the server to
  create the image. If not specified, Packer will use the environment
  variable OS_REGION_NAME, if set.

- `endpoint_type` (string) - The endpoint type to use. Can be any of "internal", "internalURL",
  "admin", "adminURL", "public", and "publicURL". By default this is
  "public".

- `cacert` (string) - Custom CA certificate file path. If omitted the OS_CACERT environment
  variable can be used.

- `cert` (string) - Client certificate file path for SSL client authentication. If omitted
  the OS_CERT environment variable can be used.

- `key` (string) - Client private key file path for SSL client authentication. If omitted
  the OS_KEY environment variable can be used.

- `token` (string) - the token (id) to use with token based authorization. Packer will use
  the environment variable OS_TOKEN, if set.

- `application_credential_name` (string) - The application credential name to use with application credential based
  authorization. Packer will use the environment variable
  OS_APPLICATION_CREDENTIAL_NAME, if set.

- `application_credential_id` (string) - The application credential id to use with application credential based
  authorization. Packer will use the environment variable
  OS_APPLICATION_CREDENTIAL_ID, if set.

- `application_credential_secret` (string) - The application credential secret to use with application credential
  based authorization. Packer will use the environment variable
  OS_APPLICATION_CREDENTIAL_SECRET, if set.

- `cloud` (string) - An entry in a `clouds.yaml` file. See the OpenStack os-client-config
  [documentation](https://docs.openstack.org/os-client-config/latest/user/configuration.html)
  for more information about `clouds.yaml` files. If omitted, the
  `OS_CLOUD` environment variable is used.

<!-- End of code generated from the comments of the AccessConfig struct in builder/openstack/access_config.go; -->


<!-- Code generated from the comments of the ImageConfig struct in builder/openstack/image_config.go; DO NOT EDIT MANUALLY -->

- `metadata` (map[string]string) - Glance metadata that will be applied to the image.

- `image_visibility` (imageservice.ImageVisibility) - One of "public", "private", "shared", or "community".

- `image_members` ([]string) - List of members to add to the image after creation. An image member is
  usually a project (also called the "tenant") with whom the image is
  shared.

- `image_auto_accept_members` (bool) - When true, perform the image accept so the members can see the image in their
  project. This requires a user with priveleges both in the build project and
  in the members provided. Defaults to false.

- `image_disk_format` (string) - Disk format of the resulting image. This option works if
  use_blockstorage_volume is true.

- `image_tags` ([]string) - List of tags to add to the image after creation.

- `image_min_disk` (int) - Minimum disk size needed to boot image, in gigabytes.

- `skip_create_image` (bool) - Skip creating the image. Useful for setting to `true` during a build test stage. Defaults to `false`.

<!-- End of code generated from the comments of the ImageConfig struct in builder/openstack/image_config.go; -->


<!-- Code generated from the comments of the RunConfig struct in builder/openstack/run_config.go; DO NOT EDIT MANUALLY -->

- `ssh_interface` (string) - The type of interface to connect via SSH. Values useful for Rackspace
  are "public" or "private", and the default behavior is to connect via
  whichever is returned first from the OpenStack API.

- `ssh_ip_version` (string) - The IP version to use for SSH connections, valid values are `4` and `6`.
  Useful on dual stacked instances where the default behavior is to
  connect via whichever IP address is returned first from the OpenStack
  API.

- `external_source_image_format` (string) - The format of the external source image to use, e.g. qcow2, raw.

- `external_source_image_properties` (map[string]string) - Properties to set for the external source image

- `availability_zone` (string) - The availability zone to launch the server in. If this isn't specified,
  the default enforced by your OpenStack cluster will be used. This may be
  required for some OpenStack clusters.

- `rackconnect_wait` (bool) - For rackspace, whether or not to wait for Rackconnect to assign the
  machine an IP address before connecting via SSH. Defaults to false.

- `floating_ip_network` (string) - The ID or name of an external network that can be used for creation of a
  new floating IP.

- `instance_floating_ip_net` (string) - The ID of the network to which the instance is attached and which should
  be used to associate with the floating IP. This provides control over
  the floating ip association on multi-homed instances. The association
  otherwise depends on a first-returned-interface policy which could fail
  if the network to which it is connected is unreachable from the floating
  IP network.

- `floating_ip` (string) - A specific floating IP to assign to this instance.

- `reuse_ips` (bool) - Whether or not to attempt to reuse existing unassigned floating ips in
  the project before allocating a new one. Note that it is not possible to
  safely do this concurrently, so if you are running multiple openstack
  builds concurrently, or if other processes are assigning and using
  floating IPs in the same openstack project while packer is running, you
  should not set this to true. Defaults to false.

- `security_groups` ([]string) - A list of security groups by name to add to this instance.

- `networks` ([]string) - A list of networks by UUID to attach to this instance.

- `ports` ([]string) - A list of ports by UUID to attach to this instance.

- `network_discovery_cidrs` ([]string) - A list of network CIDRs to discover the network to attach to this instance.
  The first network whose subnet is contained within any of the given CIDRs
  is used. Ignored if either of the above two options are provided.

- `user_data` (string) - User data to apply when launching the instance. Note that you need to be
  careful about escaping characters due to the templates being JSON. It is
  often more convenient to use user_data_file, instead. Packer will not
  automatically wait for a user script to finish before shutting down the
  instance this must be handled in a provisioner.

- `user_data_file` (string) - Path to a file that will be used for the user data when launching the
  instance.

- `instance_name` (string) - Name that is applied to the server instance created by Packer. If this
  isn't specified, the default is same as image_name.

- `instance_metadata` (map[string]string) - Metadata that is applied to the server instance created by Packer. Also
  called server properties in some documentation. The strings have a max
  size of 255 bytes each.

- `force_delete` (bool) - Whether to force the OpenStack instance to be forcefully deleted. This
  is useful for environments that have reclaim / soft deletion enabled. By
  default this is false.

- `config_drive` (bool) - Whether or not nova should use ConfigDrive for cloud-init metadata.

- `floating_ip_pool` (string) - Deprecated use floating_ip_network instead.

- `use_blockstorage_volume` (bool) - Use Block Storage service volume for the instance root volume instead of
  Compute service local volume (default).

- `volume_name` (string) - Name of the Block Storage service volume. If this isn't specified,
  random string will be used.

- `volume_type` (string) - Type of the Block Storage service volume. If this isn't specified, the
  default enforced by your OpenStack cluster will be used.

- `volume_size` (int) - Size of the Block Storage service volume in GB. If this isn't specified,
  it is set to source image min disk value (if set) or calculated from the
  source image bytes size. Note that in some cases this needs to be
  specified, if use_blockstorage_volume is true.

- `volume_availability_zone` (string) - Availability zone of the Block Storage service volume. If omitted,
  Compute instance availability zone will be used. If both of Compute
  instance and Block Storage volume availability zones aren't specified,
  the default enforced by your OpenStack cluster will be used.

- `openstack_provider` (string) - Not really used, but here for BC

- `use_floating_ip` (bool) - *Deprecated* use `floating_ip` or `floating_ip_pool` instead.

<!-- End of code generated from the comments of the RunConfig struct in builder/openstack/run_config.go; -->


### Communicator Configuration

#### Optional:

<!-- Code generated from the comments of the Config struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `communicator` (string) - Packer currently supports three kinds of communicators:
  
  -   `none` - No communicator will be used. If this is set, most
      provisioners also can't be used.
  
  -   `ssh` - An SSH connection will be established to the machine. This
      is usually the default.
  
  -   `winrm` - A WinRM connection will be established.
  
  In addition to the above, some builders have custom communicators they
  can use. For example, the Docker builder has a "docker" communicator
  that uses `docker exec` and `docker cp` to execute scripts and copy
  files.

- `pause_before_connecting` (duration string | ex: "1h5m2s") - We recommend that you enable SSH or WinRM as the very last step in your
  guest's bootstrap script, but sometimes you may have a race condition
  where you need Packer to wait before attempting to connect to your
  guest.
  
  If you end up in this situation, you can use the template option
  `pause_before_connecting`. By default, there is no pause. For example if
  you set `pause_before_connecting` to `10m` Packer will check whether it
  can connect, as normal. But once a connection attempt is successful, it
  will disconnect and then wait 10 minutes before connecting to the guest
  and beginning provisioning.

<!-- End of code generated from the comments of the Config struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSH struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `ssh_host` (string) - The address to SSH to. This usually is automatically configured by the
  builder.

- `ssh_port` (int) - The port to connect to SSH. This defaults to `22`.

- `ssh_username` (string) - The username to connect to SSH with. Required if using SSH.

- `ssh_password` (string) - A plaintext password to use to authenticate with SSH.

- `ssh_ciphers` ([]string) - This overrides the value of ciphers supported by default by Golang.
  The default value is [
    "aes128-gcm@openssh.com",
    "chacha20-poly1305@openssh.com",
    "aes128-ctr", "aes192-ctr", "aes256-ctr",
  ]
  
  Valid options for ciphers include:
  "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
  "chacha20-poly1305@openssh.com",
  "arcfour256", "arcfour128", "arcfour", "aes128-cbc", "3des-cbc",

- `ssh_clear_authorized_keys` (bool) - If true, Packer will attempt to remove its temporary key from
  `~/.ssh/authorized_keys` and `/root/.ssh/authorized_keys`. This is a
  mostly cosmetic option, since Packer will delete the temporary private
  key from the host system regardless of whether this is set to true
  (unless the user has set the `-debug` flag). Defaults to "false";
  currently only works on guests with `sed` installed.

- `ssh_key_exchange_algorithms` ([]string) - If set, Packer will override the value of key exchange (kex) algorithms
  supported by default by Golang. Acceptable values include:
  "curve25519-sha256@libssh.org", "ecdh-sha2-nistp256",
  "ecdh-sha2-nistp384", "ecdh-sha2-nistp521",
  "diffie-hellman-group14-sha1", and "diffie-hellman-group1-sha1".

- `ssh_certificate_file` (string) - Path to user certificate used to authenticate with SSH.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_pty` (bool) - If `true`, a PTY will be requested for the SSH connection. This defaults
  to `false`.

- `ssh_timeout` (duration string | ex: "1h5m2s") - The time to wait for SSH to become available. Packer uses this to
  determine when the machine has booted so this is usually quite long.
  Example value: `10m`.
  This defaults to `5m`, unless `ssh_handshake_attempts` is set.

- `ssh_disable_agent_forwarding` (bool) - If true, SSH agent forwarding will be disabled. Defaults to `false`.

- `ssh_handshake_attempts` (int) - The number of handshakes to attempt with SSH once it can connect.
  This defaults to `10`, unless a `ssh_timeout` is set.

- `ssh_bastion_host` (string) - A bastion host to use for the actual SSH connection.

- `ssh_bastion_port` (int) - The port of the bastion host. Defaults to `22`.

- `ssh_bastion_agent_auth` (bool) - If `true`, the local SSH agent will be used to authenticate with the
  bastion host. Defaults to `false`.

- `ssh_bastion_username` (string) - The username to connect to the bastion host.

- `ssh_bastion_password` (string) - The password to use to authenticate with the bastion host.

- `ssh_bastion_interactive` (bool) - If `true`, the keyboard-interactive used to authenticate with bastion host.

- `ssh_bastion_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with the
  bastion host. The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_bastion_certificate_file` (string) - Path to user certificate used to authenticate with bastion host.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_file_transfer_method` (string) - `scp` or `sftp` - How to transfer files, Secure copy (default) or SSH
  File Transfer Protocol.
  
  **NOTE**: Guests using Windows with Win32-OpenSSH v9.1.0.0p1-Beta, scp
  (the default protocol for copying data) returns a a non-zero error code since the MOTW
  cannot be set, which cause any file transfer to fail. As a workaround you can override the transfer protocol
  with SFTP instead `ssh_file_transfer_protocol = "sftp"`.

- `ssh_proxy_host` (string) - A SOCKS proxy host to use for SSH connection

- `ssh_proxy_port` (int) - A port of the SOCKS proxy. Defaults to `1080`.

- `ssh_proxy_username` (string) - The optional username to authenticate with the proxy server.

- `ssh_proxy_password` (string) - The optional password to use to authenticate with the proxy server.

- `ssh_keep_alive_interval` (duration string | ex: "1h5m2s") - How often to send "keep alive" messages to the server. Set to a negative
  value (`-1s`) to disable. Example value: `10s`. Defaults to `5s`.

- `ssh_read_write_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait for a remote command to end. This might be
  useful if, for example, packer hangs on a connection after a reboot.
  Example: `5m`. Disabled by default.

- `ssh_remote_tunnels` ([]string) - 

- `ssh_local_tunnels` ([]string) - 

<!-- End of code generated from the comments of the SSH struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `temporary_key_pair_type` (string) - `dsa` | `ecdsa` | `ed25519` | `rsa` ( the default )
  
  Specifies the type of key to create. The possible values are 'dsa',
  'ecdsa', 'ed25519', or 'rsa'.
  
  NOTE: DSA is deprecated and no longer recognized as secure, please
  consider other alternatives like RSA or ED25519.

- `temporary_key_pair_bits` (int) - Specifies the number of bits in the key to create. For RSA keys, the
  minimum size is 1024 bits and the default is 4096 bits. Generally, 3072
  bits is considered sufficient. DSA keys must be exactly 1024 bits as
  specified by FIPS 186-2. For ECDSA keys, bits determines the key length
  by selecting from one of three elliptic curve sizes: 256, 384 or 521
  bits. Attempting to use bit lengths other than these three values for
  ECDSA keys will fail. Ed25519 keys have a fixed length and bits will be
  ignored.
  
  NOTE: DSA is deprecated and no longer recognized as secure as specified
  by FIPS 186-5, please consider other alternatives like RSA or ED25519.

<!-- End of code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; -->


- `ssh_keypair_name` (string) - If specified, this is the key that will be used for SSH with the
  machine. The key must match a key pair name loaded up into the remote.
  By default, this is blank, and Packer will generate a temporary keypair
  unless [`ssh_password`](#ssh_password) is used.
  [`ssh_private_key_file`](#ssh_private_key_file) or
  [`ssh_agent_auth`](#ssh_agent_auth) must be specified when
  [`ssh_keypair_name`](#ssh_keypair_name) is utilized.


- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


- `ssh_agent_auth` (bool) - If true, the local SSH agent will be used to authenticate connections to
  the source instance. No temporary keypair will be created, and the
  values of [`ssh_password`](#ssh_password) and
  [`ssh_private_key_file`](#ssh_private_key_file) will be ignored. The
  environment variable `SSH_AUTH_SOCK` must be set for this option to work
  properly.


<!-- Code generated from the comments of the SSHInterface struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `ssh_interface` (string) - One of `public_ip`, `private_ip`, `public_dns`, or `private_dns`. If
  set, either the public IP address, private IP address, public DNS name
  or private DNS name will used as the host for SSH. The default behaviour
  if inside a VPC is to use the public IP address if available, otherwise
  the private IP address will be used. If not in a VPC the public DNS name
  will be used. Also works for WinRM.
  
  Where Packer is configured for an outbound proxy but WinRM traffic
  should be direct, `ssh_interface` must be set to `private_dns` and
  `<region>.compute.internal` included in the `NO_PROXY` environment
  variable.

- `ssh_ip_version` (string) - The IP version to use for SSH connections, valid values are `4` and `6`.
  Useful on dual stacked instances where the default behavior is to
  connect via whichever IP address is returned first from the OpenStack
  API.

<!-- End of code generated from the comments of the SSHInterface struct in communicator/config.go; -->


## Basic Example: DevStack

Here is a basic example. This is a example to build on DevStack running in a
VM.

**JSON**

```json
{
  "builders":
  [{
    "type": "openstack",
    "identity_endpoint": "http://<devstack-ip>:5000/v3",
    "tenant_name": "admin",
    "domain_name": "Default",
    "username": "admin",
    "password": "<your admin password>",
    "region": "RegionOne",
    "ssh_username": "root",
    "image_name": "Test image",
    "source_image": "<image id>",
    "flavor": "m1.tiny",
    "insecure": "true"
  }]
}
```

**HCL2**

```hcl
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

```


## Basic Example: Rackspace public cloud

Here is a basic example. This is a working example to build a Ubuntu 12.04 LTS
(Precise Pangolin) on Rackspace OpenStack cloud offering.

**JSON**

```json
{
  "builders":[{
    "type": "openstack",
    "username": "foo",
    "password": "foo",
    "region": "DFW",
    "ssh_username": "root",
    "image_name": "Test image",
    "source_image": "23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
    "flavor": "2"
  }]
}
```


**HCL2**

```hcl

source "openstack" "example" {
  flavor       = "2"
  image_name   = "Test image"
  password     = "foo"
  region       = "DFW"
  source_image = "23b564c9-c3e6-49f9-bc68-86c7a9ab5018"
  ssh_username = "root"
  username     = "foo"
}

build {
  sources = ["source.openstack.example"]
}
```


## Basic Example: Private OpenStack cloud

This example builds an Ubuntu 14.04 image on a private OpenStack cloud, powered
by Metacloud.

**JSON**

```json
{
  "builders":[{
    "type": "openstack",
    "ssh_username": "root",
    "image_name": "ubuntu1404_packer_test_1",
    "source_image": "91d9c168-d1e5-49ca-a775-3bfdbb6c97f1",
    "flavor": "2"
  }]
}
```


**HCL2**

```hcl
source "openstack" "example" {
  flavor       = "2"
  image_name   = "ubuntu1404_packer_test_1"
  source_image = "91d9c168-d1e5-49ca-a775-3bfdbb6c97f1"
  ssh_username = "root"
}

build {
  sources = ["source.openstack.example"]
}
```


In this case, the connection information for connecting to OpenStack doesn't
appear in the template. That is because I source a standard OpenStack script
with environment variables set before I run this. This script is setting
environment variables like:

- `OS_AUTH_URL`
- `OS_TENANT_ID`
- `OS_USERNAME`
- `OS_PASSWORD`

This is slightly different when identity v3 is used:

- `OS_AUTH_URL`
- `OS_USERNAME`
- `OS_PASSWORD`
- `OS_DOMAIN_NAME`
- `OS_TENANT_NAME`

This will authenticate the user on the domain and scope you to the project. A
tenant is the same as a project. It's optional to use names or IDs in v3. This
means you can use `OS_USERNAME` or `OS_USERID`, `OS_TENANT_ID` or
`OS_TENANT_NAME` and `OS_DOMAIN_ID` or `OS_DOMAIN_NAME`.

The above example would be equivalent to an RC file looking like this :

```shell
export OS_AUTH_URL="https://identity.myprovider/v3"
export OS_USERNAME="myuser"
export OS_PASSWORD="password"
export OS_USER_DOMAIN_NAME="mydomain"
export OS_PROJECT_DOMAIN_NAME="mydomain"
```

## Basic Example: Instance with Block Storage root volume

A basic example of Instance with a remote root Block Storage service volume.
This is a working example to build an image on private OpenStack cloud powered
by Selectel VPC.

**JSON**

```json
{
  "builders":[{
    "type": "openstack",
    "identity_endpoint": "https://api.selvpc.com/identity/v3",
    "tenant_id": "2e90c5c04c7b4c509be78723e2b55b77",
    "username": "foo",
    "password": "foo",
    "region": "ru-3",
    "ssh_username": "root",
    "image_name": "Test image",
    "source_image": "5f58ea7e-6264-4939-9d0f-0c23072b1132",
    "networks": "9aab504e-bedf-48af-9256-682a7fa3dabb",
    "flavor": "1001",
    "availability_zone": "ru-3a",
    "use_blockstorage_volume": true,
    "volume_type": "fast.ru-3a"
  }]
}
```


**HCL2**

```hcl

source "openstack" "example" {
  availability_zone       = "ru-3a"
  flavor                  = "1001"
  identity_endpoint       = "https://api.selvpc.com/identity/v3"
  image_name              = "Test image"
  networks                = "9aab504e-bedf-48af-9256-682a7fa3dabb"
  password                = "foo"
  region                  = "ru-3"
  source_image            = "5f58ea7e-6264-4939-9d0f-0c23072b1132"
  ssh_username            = "root"
  tenant_id               = "2e90c5c04c7b4c509be78723e2b55b77"
  use_blockstorage_volume = true
  username                = "foo"
  volume_type             = "fast.ru-3a"
}

build {
  sources = ["source.openstack.example"]
}

```


## Notes on OpenStack Authorization

The simplest way to get all settings for authorization against OpenStack is to
go into the OpenStack Dashboard (Horizon) select your _Project_ and navigate
_Project, Access & Security_, select _API Access_ and _Download OpenStack RC
File v3_. Source the file, and select your wanted region by setting environment
variable `OS_REGION_NAME` or `OS_REGION_ID` and
`export OS_TENANT_NAME=$OS_PROJECT_NAME` or
`export OS_TENANT_ID=$OS_PROJECT_ID`.

~> `OS_TENANT_NAME` or `OS_TENANT_ID` must be used even with Identity v3,
`OS_PROJECT_NAME` and `OS_PROJECT_ID` has no effect in Packer.

To troubleshoot authorization issues test you environment variables with the
OpenStack cli. It can be installed with

    $ pip install --user python-openstackclient

### Authorize Using Tokens

To authorize with a access token only `identity_endpoint` and `token` is
needed, and possibly `tenant_name` or `tenant_id` depending on your token type.
Or use the following environment variables:

- `OS_AUTH_URL`
- `OS_TOKEN`
- One of `OS_TENANT_NAME` or `OS_TENANT_ID`

### Authorize Using Application Credential

To authorize with an application credential, only `identity_endpoint`,
`application_credential_id`, and `application_credential_secret` are needed.
Or use the following environment variables:

- `OS_AUTH_URL`
- `OS_APPLICATION_CREDENTIAL_ID`
- `OS_APPLICATION_CREDENTIAL_SECRET`
