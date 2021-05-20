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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestFirewall(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	firewallName := config.AppendRandomSuffix("firewall")
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	firewallPolicyName := config.AppendRandomSuffix("firewallpolicy")

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWANPro := armnetwork.VirtualWAN{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualWanProperties{
			DisableVPNEncryption: to.BoolPtr(false),
			Type:                 to.StringPtr("Basic"),
		},
	}
	virtualWanID, err := CreateVirtualWan(ctx, virtualWanName, virtualWANPro)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}
	t.Logf("created virtual wan")

	virtualHubPro := armnetwork.VirtualHub{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualHubProperties{
			AddressPrefix: to.StringPtr("10.168.0.0/24"),
			SKU:           to.StringPtr("Basic"),
			VirtualWan: &armnetwork.SubResource{
				ID: &virtualWanID,
			},
		},
	}
	virtualHubId, err := CreateVirtualHub(ctx, virtualHubName, virtualWanID, virtualHubPro)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}
	t.Logf("created virtual hub")

	firewallPolicyPro := armnetwork.FirewallPolicy{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.FirewallPolicyPropertiesFormat{
			ThreatIntelMode: armnetwork.AzureFirewallThreatIntelModeAlert.ToPtr(),
		},
	}
	firewallPolicyId, err := CreateFirewallPolicy(ctx, firewallPolicyName, firewallPolicyPro)
	if err != nil {
		t.Fatalf("failed to create firewall policy: % +v", err)
	}
	t.Logf("created firewall policy")

	azureFirewallPro := armnetwork.AzureFirewall{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.AzureFirewallPropertiesFormat{
			SKU: &armnetwork.AzureFirewallSKU{
				Name: armnetwork.AzureFirewallSKUNameAZFWHub.ToPtr(),
				Tier: armnetwork.AzureFirewallSKUTierStandard.ToPtr(),
			},
			VirtualHub: &armnetwork.SubResource{
				ID: &virtualHubId,
			},
			FirewallPolicy: &armnetwork.SubResource{
				ID: &firewallPolicyId,
			},
			HubIPAddresses: &armnetwork.HubIPAddresses{
				PublicIPs: &armnetwork.HubPublicIPAddresses{
					Addresses: &[]*armnetwork.AzureFirewallPublicIPAddress{},
					Count:     to.Int32Ptr(1),
				},
			},
		},
		Zones: &[]*string{},
	}
	err = CreateFirewall(ctx, firewallName, azureFirewallPro)
	if err != nil {
		t.Fatalf("failed to create firewall: % +v", err)
	}
	t.Logf("created firewall")

	err = GetFirewall(ctx, firewallName)
	if err != nil {
		t.Fatalf("failed to get firewall: %+v", err)
	}
	t.Logf("got firewall")

	err = ListFirewall(ctx)
	if err != nil {
		t.Fatalf("failed to list firewall: %+v", err)
	}
	t.Logf("listed firewall")

	err = ListAllAzureFirewallFqdnTag(ctx)
	if err != nil {
		t.Fatalf("failed to list all azure firewall fqdn tag: %+v", err)
	}
	t.Logf("listed all azure firewall fqdn tag")

	err = ListAllFirewall(ctx)
	if err != nil {
		t.Fatalf("failed to list all firewall: %+v", err)
	}
	t.Logf("listed all firewall")
	err = UpdateFirewallTags(ctx, firewallName)
	if err != nil {
		t.Fatalf("failed to update tags for firewall: %+v", err)
	}
	t.Logf("updated firewall tags")

	err = DeleteFirewall(ctx, firewallName)
	if err != nil {
		t.Fatalf("failed to delete firewall: %+v", err)
	}
	t.Logf("deleted firewall")

}
