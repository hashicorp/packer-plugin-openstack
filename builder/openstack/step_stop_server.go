// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepStopServer struct{}

func (s *StepStopServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)
	server := state.Get("server").(*servers.Server)

	// We need the v2 compute client
	client, err := config.computeV2Client()
	if err != nil {
		err = fmt.Errorf("Error initializing compute client: %s", err)
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Stopping server: %s ...", server.ID))
	if err := servers.Stop(ctx, client, server.ID).ExtractErr(); err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
			// The server might have already been shut down by Windows Sysprep
			log.Printf("[WARN] 409 on stopping an already stopped server, continuing")
			return multistep.ActionContinue
		} else {
			err = fmt.Errorf("Error stopping server: %s", err)
			state.Put("error", err)
			return multistep.ActionHalt
		}
	}

	ui.Say(fmt.Sprintf("Waiting for server to stop: %s ...", server.ID))
	stateChange := StateChangeConf{
		Pending:   []string{"ACTIVE"},
		Target:    []string{"SHUTOFF", "STOPPED"},
		Refresh:   ServerStateRefreshFunc(ctx, client, server.ID),
		StepState: state,
	}
	if _, err := WaitForState(&stateChange); err != nil {
		err := fmt.Errorf("Error waiting for server (%s) to stop: %s", server.ID, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *StepStopServer) Cleanup(state multistep.StateBag) {}
