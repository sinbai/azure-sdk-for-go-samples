// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

func TestFirewallPolicy(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	firewallPolicyName := config.AppendRandomSuffix("firewallpolicy")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to create firewall policy: % +v", err)
	}
	t.Logf("created firewall policy")

	err = GetFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to get firewall policy: %+v", err)
	}
	t.Logf("got firewall policy")

	err = ListFirewallPolicy(ctx)
	if err != nil {
		t.Fatalf("failed to list firewall policy: %+v", err)
	}
	t.Logf("listed firewall policy")

	err = ListAllFirewallPolicy(ctx)
	if err != nil {
		t.Fatalf("failed to list all firewall policy: %+v", err)
	}
	t.Logf("listed all firewall policy")

	err = DeleteFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to delete firewall policy: %+v", err)
	}
	t.Logf("deleted firewall policy")

}
