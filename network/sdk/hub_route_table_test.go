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

func TestHubRouteTable(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	hubRouteTableName := config.AppendRandomSuffix("hubroutetable")
	firewallName := config.AppendRandomSuffix("firewall")
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	firewallPolicyName := config.AppendRandomSuffix("firewallpolicy")

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWANParameters := armnetwork.VirtualWAN{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualWanProperties{
			DisableVPNEncryption: to.BoolPtr(false),
			Type:                 to.StringPtr("Basic"),
		},
	}
	virtualWanID, err := CreateVirtualWan(ctx, virtualWanName, virtualWANParameters)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}

	virtualHubParameters := armnetwork.VirtualHub{
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
	virtualHubId, err := CreateVirtualHub(ctx, virtualHubName, virtualHubParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}

	firewallPolicyParameters := armnetwork.FirewallPolicy{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.FirewallPolicyPropertiesFormat{
			ThreatIntelMode: armnetwork.AzureFirewallThreatIntelModeAlert.ToPtr(),
		},
	}
	firewallPolicyId, err := CreateFirewallPolicy(ctx, firewallPolicyName, firewallPolicyParameters)
	if err != nil {
		t.Fatalf("failed to create firewall policy: % +v", err)
	}

	azureFirewallParameters := armnetwork.AzureFirewall{
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
					Count: to.Int32Ptr(1),
				},
			},
		},
	}
	firewallId, err := CreateFirewall(ctx, firewallName, azureFirewallParameters)
	if err != nil {
		t.Fatalf("failed to create firewall: % +v", err)
	}

	routeTableParameters := armnetwork.HubRouteTable{
		Properties: &armnetwork.HubRouteTableProperties{
			Labels: &[]*string{to.StringPtr("label1"),
				to.StringPtr("label2")},
			Routes: &[]*armnetwork.HubRoute{{
				Name:            to.StringPtr("route1"),
				DestinationType: to.StringPtr("CIDR"),
				Destinations: &[]*string{
					to.StringPtr("10.0.0.0/8"),
					to.StringPtr("20.0.0.0/8"),
					to.StringPtr("30.0.0.0/8")},
				NextHop:     &firewallId,
				NextHopType: to.StringPtr("ResourceId"),
			}},
		},
	}
	err = CreateHubRouteTable(ctx, virtualHubName, hubRouteTableName, routeTableParameters)
	if err != nil {
		t.Fatalf("failed to create hub route table: % +v", err)
	}
	t.Logf("created hub route table")

	err = GetHubRouteTable(ctx, virtualHubName, hubRouteTableName)
	if err != nil {
		t.Fatalf("failed to get hub route table: %+v", err)
	}
	t.Logf("got hub route table")

	err = ListHubRouteTable(ctx, virtualHubName)
	if err != nil {
		t.Fatalf("failed to list hub route table: %+v", err)
	}
	t.Logf("listed hub route table")

	err = DeleteHubRouteTable(ctx, virtualHubName, hubRouteTableName)
	if err != nil {
		t.Fatalf("failed to delete hub route table: %+v", err)
	}
	t.Logf("deleted hub route table")
}
