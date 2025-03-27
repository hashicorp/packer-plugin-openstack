// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepDeleteServer struct {
	UseBlockStorageVolume bool
}

func (s *StepDeleteServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// Only if we have a blockstorage volume we need to detach it to upload it as an image
	if !s.UseBlockStorageVolume {
		return multistep.ActionContinue
	}

	instance := state.Get("instance_id").(string)

	err := DeleteServer(state, instance)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	state.Put("instance_id", "")
	return multistep.ActionContinue
}

func (s *StepDeleteServer) Cleanup(state multistep.StateBag) {
}
