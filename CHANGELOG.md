## 1.1.4 (June 23, 2026)

* fix: add an optional wait before creating the image to improve reliability in image creation workflows.
* fix: use `openstack` as the HCP metadata provider name (instead of the default prefixed provider value).
* chore: add the HCP Ready flag to web metadata.
* chore: update copyright and license headers across the codebase for compliance.
* chore: bump `github.com/hashicorp/packer-plugin-sdk` from `0.6.4` to `0.6.7`.

## 1.0.0 (June 14, 2021)

* Update packer-plugin-sdk to version 0.2.3. [GH-29]

## 0.0.2 (April 21, 2021)

* fast-follow release to fix goreleaser on github.

## 0.0.1 (April 21, 2021)

* Hyper-V Plugin break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md).