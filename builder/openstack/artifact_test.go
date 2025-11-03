// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"reflect"
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
	"github.com/mitchellh/mapstructure"
)

func TestArtifact_Impl(t *testing.T) {
	var _ packersdk.Artifact = new(Artifact)
}

func TestArtifactId(t *testing.T) {
	expected := `b8cdf55b-c916-40bd-b190-389ec144c4ed`

	a := &Artifact{
		ImageId: "b8cdf55b-c916-40bd-b190-389ec144c4ed",
	}

	result := a.Id()
	if result != expected {
		t.Fatalf("bad: %s", result)
	}
}

func TestArtifactString(t *testing.T) {
	expected := "An image was created: b8cdf55b-c916-40bd-b190-389ec144c4ed"

	a := &Artifact{
		ImageId: "b8cdf55b-c916-40bd-b190-389ec144c4ed",
	}
	result := a.String()
	if result != expected {
		t.Fatalf("bad: %s", result)
	}
}

func TestArtifactState_StateData(t *testing.T) {
	expectedData := "this is the data"
	artifact := &Artifact{
		StateData: map[string]interface{}{"state_data": expectedData},
	}

	// Valid state
	result := artifact.State("state_data")
	if result != expectedData {
		t.Fatalf("Bad: State data was %s instead of %s", result, expectedData)
	}

	// Invalid state
	result = artifact.State("invalid_key")
	if result != nil {
		t.Fatalf("Bad: State should be nil for invalid state data name")
	}

	// Nil StateData should not fail and should return nil
	artifact = &Artifact{}
	result = artifact.State("key")
	if result != nil {
		t.Fatalf("Bad: State should be nil for nil StateData")
	}
}

func TestArtifactState_hcpPackerRegistryMetadata(t *testing.T) {
	artifact := &Artifact{
		ImageId:        "foo",
		BuilderIdValue: "openstack",
		SourceImage:    "bar",
		Region:         "mordor-7",
	}

	result := artifact.State(registryimage.ArtifactStateURI)
	if result == nil {
		t.Fatalf("Bad: HCP Packer registry image data was nil")
	}

	var image registryimage.Image
	err := mapstructure.Decode(result, &image)
	if err != nil {
		t.Errorf("Bad: unexpected error when trying to decode state into registryimage.Image %v", err)
	}

	expected := registryimage.Image{
		ImageID:        "foo",
		ProviderName:   "openstack",
		ProviderRegion: "mordor-7",
		SourceImageID:  "bar",
		Labels:         map[string]string{},
	}

	if !reflect.DeepEqual(image, expected) {
		t.Fatalf("bad: \ngot:\n%#v,\nexpected: \n%#v", image, expected)
	}
}
