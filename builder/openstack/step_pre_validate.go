// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepPreValidate struct {
	ForceImageName    bool
}

func (s *StepPreValidate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	client, err := config.imageV2Client()
	if err != nil {
		err := fmt.Errorf("error creating image client: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if s.ForceImageName {
		ui.Say("ForceImageName flag found, skipping prevalidating Image Name")
		return multistep.ActionContinue
	}

	ui.Say(fmt.Sprintf("Prevalidating Image Name: %s", config.ImageName))

	listOpts := images.ListOpts{
		Name: config.ImageName,
	}	

	var imageList []images.Image
	err = images.List(client, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		imageBatch, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}
		imageList = append(imageList, imageBatch...)
		return true, nil
	})

	if err != nil {
		err := fmt.Errorf("Error querying image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(imageList) > 0 {
		err := fmt.Errorf("Error: Image Name: '%s' has already been used", config.ImageName)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt		
	}

	return multistep.ActionContinue
}

func (s *StepPreValidate) Cleanup(state multistep.StateBag) {
}
