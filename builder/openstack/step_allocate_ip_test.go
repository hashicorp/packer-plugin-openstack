// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

const testFloatingNetworkID = "9a8b7c6d-0000-4000-8000-000000000001"

type allocateIpCalls struct {
	listed, created, deleted int
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func allocateIpFixture(t *testing.T, listBody string, listStatus int) *allocateIpCalls {
	th.SetupHTTP()
	t.Cleanup(th.TeardownHTTP)
	calls := &allocateIpCalls{}

	th.Mux.HandleFunc("GET /v2.0/floatingips", func(w http.ResponseWriter, r *http.Request) {
		calls.listed++
		writeJSON(w, listStatus, `{"floatingips": `+listBody+`}`)
	})
	th.Mux.HandleFunc("POST /v2.0/floatingips", func(w http.ResponseWriter, r *http.Request) {
		calls.created++
		th.TestJSONRequest(t, r, fmt.Sprintf(`{"floatingip": {"floating_network_id": %q}}`, testFloatingNetworkID))
		writeJSON(w, http.StatusCreated, `{"floatingip": {"id": "fip-new", "floating_ip_address": "203.0.113.10", "status": "DOWN"}}`)
	})
	th.Mux.HandleFunc("PUT /v2.0/floatingips/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, fmt.Sprintf(`{"floatingip": {"id": %q, "port_id": "port-1"}}`, r.PathValue("id")))
	})
	th.Mux.HandleFunc("DELETE /v2.0/floatingips/{id}", func(w http.ResponseWriter, r *http.Request) {
		calls.deleted++
		w.WriteHeader(http.StatusNoContent)
	})
	th.Mux.HandleFunc("GET /servers/{id}/os-interface", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{"interfaceAttachments": [{"port_id": "port-1", "net_id": "net-1", "port_state": "ACTIVE", "mac_addr": "fa:16:3e:00:00:01", "fixed_ips": []}]}`)
	})

	return calls
}

func allocateIpState(t *testing.T) multistep.StateBag {
	config := &Config{}
	config.osClient = &gophercloud.ProviderClient{
		TokenID: "token",
		EndpointLocator: func(gophercloud.EndpointOpts) (string, error) {
			return th.Endpoint(), nil
		},
	}

	state := new(multistep.BasicStateBag)
	state.Put("ui", packersdk.TestUi(t))
	state.Put("config", config)
	state.Put("server", &servers.Server{ID: "server-1"})
	return state
}

func TestStepAllocateIp_ReuseIPsAllocatesWhenNoneFree(t *testing.T) {
	calls := allocateIpFixture(t, `[]`, http.StatusOK)
	state := allocateIpState(t)
	step := &StepAllocateIp{ReuseIPs: true, FloatingIPNetwork: testFloatingNetworkID}

	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %v, error: %v", action, state.Get("error"))
	}
	if ip := state.Get("access_ip").(*floatingips.FloatingIP); ip.ID != "fip-new" {
		t.Fatalf("expected the allocated floating IP, got %q", ip.ID)
	}
	if state.Get("floatingip_istemp") != true {
		t.Fatalf("allocated floating IP should be marked temporary")
	}
	if calls.listed != 1 || calls.created != 1 {
		t.Fatalf("expected one list and one create call, got %d and %d", calls.listed, calls.created)
	}

	step.Cleanup(state)
	if calls.deleted != 1 {
		t.Fatalf("expected the allocated floating IP to be deleted, got %d delete calls", calls.deleted)
	}
}

func TestStepAllocateIp_ReuseIPsUsesFreeIP(t *testing.T) {
	calls := allocateIpFixture(t, `[{"id": "fip-free", "floating_ip_address": "203.0.113.5", "port_id": "", "status": "DOWN"}]`, http.StatusOK)
	state := allocateIpState(t)
	step := &StepAllocateIp{ReuseIPs: true, FloatingIPNetwork: testFloatingNetworkID}

	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %v, error: %v", action, state.Get("error"))
	}
	if ip := state.Get("access_ip").(*floatingips.FloatingIP); ip.ID != "fip-free" {
		t.Fatalf("expected the free floating IP, got %q", ip.ID)
	}
	if state.Get("floatingip_istemp") != false {
		t.Fatalf("reused floating IP should not be marked temporary")
	}
	if calls.created != 0 {
		t.Fatalf("expected no create call, got %d", calls.created)
	}

	step.Cleanup(state)
	if calls.deleted != 0 {
		t.Fatalf("reused floating IP should not be deleted, got %d delete calls", calls.deleted)
	}
}

func TestStepAllocateIp_ReuseIPsNoneFreeWithoutNetwork(t *testing.T) {
	calls := allocateIpFixture(t, `[]`, http.StatusOK)
	state := allocateIpState(t)
	step := &StepAllocateIp{ReuseIPs: true}

	if action := step.Run(context.Background(), state); action != multistep.ActionHalt {
		t.Fatalf("bad action: %v", action)
	}
	err, _ := state.Get("error").(error)
	if err == nil || !strings.Contains(err.Error(), "no free floating IPs found") || !strings.Contains(err.Error(), "floating_ip_network") {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls.created != 0 {
		t.Fatalf("expected no create call, got %d", calls.created)
	}

	step.Cleanup(state)
	if calls.deleted != 0 {
		t.Fatalf("expected no delete call, got %d", calls.deleted)
	}
}

func TestStepAllocateIp_ReuseIPsListErrorHalts(t *testing.T) {
	calls := allocateIpFixture(t, `[]`, http.StatusInternalServerError)
	state := allocateIpState(t)
	step := &StepAllocateIp{ReuseIPs: true, FloatingIPNetwork: testFloatingNetworkID}

	if action := step.Run(context.Background(), state); action != multistep.ActionHalt {
		t.Fatalf("bad action: %v", action)
	}
	if state.Get("error") == nil {
		t.Fatalf("expected an error in the state bag")
	}
	if calls.created != 0 {
		t.Fatalf("expected no create call after a list error, got %d", calls.created)
	}
}

func TestStepAllocateIp_FloatingIPNetworkWithoutReuse(t *testing.T) {
	calls := allocateIpFixture(t, `[{"id": "fip-free", "floating_ip_address": "203.0.113.5", "port_id": "", "status": "DOWN"}]`, http.StatusOK)
	state := allocateIpState(t)
	step := &StepAllocateIp{FloatingIPNetwork: testFloatingNetworkID}

	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %v, error: %v", action, state.Get("error"))
	}
	if ip := state.Get("access_ip").(*floatingips.FloatingIP); ip.ID != "fip-new" {
		t.Fatalf("expected the allocated floating IP, got %q", ip.ID)
	}
	if calls.listed != 0 || calls.created != 1 {
		t.Fatalf("expected no list and one create call, got %d and %d", calls.listed, calls.created)
	}
}
