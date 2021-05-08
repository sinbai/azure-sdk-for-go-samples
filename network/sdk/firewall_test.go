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

func TestFirewall(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	config.SetLocation("West US")
	firewallName := config.AppendRandomSuffix("firewall")
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	firewallPolicyName := config.AppendRandomSuffix("firewallpolicy")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.DefaultLocation())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	_, err = CreateVirtualWan(ctx, virtualWanName)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}
	t.Logf("created virtual wan")

	err = CreateVirtualHub(ctx, virtualHubName, virtualWanName)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}
	t.Logf("created virtual hub")

	err = CreateFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to create firewall policy: % +v", err)
	}
	t.Logf("created firewall policy")

	err = CreateFirewall(ctx, firewallName, firewallPolicyName, virtualHubName)
	if err != nil {
		t.Fatalf("failed to create firewall: % +v", err)
	}
	t.Logf("created firewall")

	/********the following code need to test after firewall created successfully*********/
	// err = GetFirewall(ctx, firewallName)
	// if err != nil {
	// 	t.Fatalf("failed to get firewall: %+v", err)
	// }
	// t.Logf("got firewall")

	// err = ListFirewall(ctx)
	// if err != nil {
	// 	t.Fatalf("failed to list firewall: %+v", err)
	// }
	// t.Logf("listed firewall")

	// err = ListAllAzureFirewallFqdnTag(ctx)
	// if err != nil {
	// 	t.Fatalf("failed to list all azure firewall fqdn tag: %+v", err)
	// }
	// t.Logf("listed all azure firewall fqdn tag")

	// err = ListAllFirewall(ctx)
	// if err != nil {
	// 	t.Fatalf("failed to list all firewall: %+v", err)
	// }
	// t.Logf("listed all firewall")
	// err = UpdateFirewallTags(ctx, firewallName)
	// if err != nil {
	// 	t.Fatalf("failed to update tags for firewall: %+v", err)
	// }
	// t.Logf("updated firewall tags")

	// err = DeleteFirewall(ctx, firewallName)
	// if err != nil {
	// 	t.Fatalf("failed to delete firewall: %+v", err)
	// }
	// t.Logf("deleted firewall")

}
