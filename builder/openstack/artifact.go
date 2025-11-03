// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
)

// Artifact is an artifact implementation that contains built images.
type Artifact struct {
	// ImageId of built image
	ImageId string

	// BuilderId is the unique ID for the builder that created this image
	BuilderIdValue string

	// OpenStack connection for performing API stuff.
	Client *gophercloud.ServiceClient

	// SourceImage ID of the created image this is actually resolved
	// based on a few configured config attributes
	SourceImage string

	// The region is read from the env OS_REGION_NAME if not provided
	Region string

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]any
}

func (a *Artifact) BuilderId() string {
	return a.BuilderIdValue
}

func (*Artifact) Files() []string {
	// We have no files
	return nil
}

func (a *Artifact) Id() string {
	return a.ImageId
}

func (a *Artifact) String() string {
	return fmt.Sprintf("An image was created: %v", a.ImageId)
}

func (a *Artifact) State(name string) any {

	if name == registryimage.ArtifactStateURI {
		img, err := registryimage.FromArtifact(a,
			registryimage.WithRegion(a.Region),
			registryimage.WithSourceID(a.SourceImage),
		)

		if err != nil {
			log.Printf("[DEBUG] error encountered when creating a registry image %v", err)
			return nil
		}
		return img

	}
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %s", a.ImageId)
	return images.Delete(a.Client, a.ImageId).ExtractErr()
}
