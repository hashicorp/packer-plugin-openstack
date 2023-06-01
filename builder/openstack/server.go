// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// StateRefreshFunc is a function type used for StateChangeConf that is
// responsible for refreshing the item being watched for a state change.
//
// It returns three results. `result` is any object that will be returned
// as the final object after waiting for state change. This allows you to
// return the final updated object, for example an openstack instance after
// refreshing it.
//
// `state` is the latest state of that object. And `err` is any error that
// may have happened while refreshing the state.
type StateRefreshFunc func() (result interface{}, state string, progress int, err error)

// StateChangeConf is the configuration struct used for `WaitForState`.
type StateChangeConf struct {
	Pending   []string
	Refresh   StateRefreshFunc
	StepState multistep.StateBag
	Target    []string
}

// ServerStateRefreshFunc returns a StateRefreshFunc that is used to watch
// an openstack server.
func ServerStateRefreshFunc(
	client *gophercloud.ServiceClient, s *servers.Server) StateRefreshFunc {
	return func() (interface{}, string, int, error) {
		serverNew, err := servers.Get(client, s.ID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[INFO] 404 on ServerStateRefresh, returning DELETED")
				return nil, "DELETED", 0, nil
			}
			log.Printf("[ERROR] Error on ServerStateRefresh: %s", err)
			return nil, "", 0, err
		}

		return serverNew, serverNew.Status, serverNew.Progress, nil
	}
}

// WaitForState watches an object and waits for it to achieve a certain
// state.
func WaitForState(conf *StateChangeConf) (i interface{}, err error) {
	log.Printf("Waiting for state to become: %s", conf.Target)

	for {
		var currentProgress int
		var currentState string
		i, currentState, currentProgress, err = conf.Refresh()
		if err != nil {
			return
		}

		for _, t := range conf.Target {
			if currentState == t {
				return
			}
		}

		if conf.StepState != nil {
			if _, ok := conf.StepState.GetOk(multistep.StateCancelled); ok {
				return nil, errors.New("interrupted")
			}
		}

		found := false
		for _, allowed := range conf.Pending {
			if currentState == allowed {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("unexpected state '%s', wanted target '%s'", currentState, conf.Target)
		}

		log.Printf("Waiting for state to become: %s currently %s (%d%%)", conf.Target, currentState, currentProgress)
		time.Sleep(2 * time.Second)
	}
}

func DeleteServer(state multistep.StateBag, instance string) error {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	// We need the v2 compute client
	computeClient, err := config.computeV2Client()
	if err != nil {
		err = fmt.Errorf("Error terminating server, may still be around: %s", err)
		return err
	}

	maxNumErrors := 10
	numErrors := 0

	ui.Say(fmt.Sprintf("Terminating the source server: %s ...", instance))
	for {
		if config.ForceDelete {
			err = servers.ForceDelete(computeClient, instance).ExtractErr()
		} else {
			err = servers.Delete(computeClient, instance).ExtractErr()
		}

		if err == nil {
			break
		}

		if _, ok := err.(gophercloud.ErrDefault500); !ok {
			err = fmt.Errorf("Error terminating server, may still be around: %s", err)
			return err
		}

		if numErrors < maxNumErrors {
			numErrors++
			log.Printf("Error terminating server on (%d) time(s): %s, retrying ...", numErrors, err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = fmt.Errorf("Error terminating server, maximum number (%d) reached: %s", numErrors, err)
		return err
	}

	server, err := servers.Get(computeClient, instance).Extract()
	if err != nil {
		err = fmt.Errorf("Error getting server to terminate: %s", err)
		return err
	}

	stateChange := StateChangeConf{
		Pending: []string{"ACTIVE", "BUILD", "REBUILD", "SUSPENDED", "SHUTOFF", "STOPPED"},
		Refresh: ServerStateRefreshFunc(computeClient, server),
		Target:  []string{"DELETED"},
	}

	_, err = WaitForState(&stateChange)
	if err != nil {
		err = fmt.Errorf("Error terminating server: %s", err)
		return err
	}
	return nil
}
