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
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVirtualNetworkPeering(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualNetworkPeeringName := config.AppendRandomSuffix("virtualnetworkpeering")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	remoteVirtualNetworkName := config.AppendRandomSuffix("remotevirtualnetwork")
	subNetName := config.AppendRandomSuffix("subnet")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualNetworkParameters := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{to.StringPtr("10.2.0.0/16")},
			},
		},
	}
	remoteVirtualNetworkId, err := CreateVirtualNetwork(ctx, remoteVirtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	virtualNetworkParameters = armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                     to.StringPtr("10.0.1.0/24"),
			PrivateLinkServiceNetworkPolicies: armnetwork.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled.ToPtr(),
		},
	}
	_, err = CreateSubnet(ctx, virtualNetworkName, subNetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	virtualNetworkPeeringParameters := armnetwork.VirtualNetworkPeering{
		Properties: &armnetwork.VirtualNetworkPeeringPropertiesFormat{
			AllowForwardedTraffic:     to.BoolPtr(true),
			AllowGatewayTransit:       to.BoolPtr(false),
			AllowVirtualNetworkAccess: to.BoolPtr(true),
			RemoteVirtualNetwork: &armnetwork.SubResource{
				ID: &remoteVirtualNetworkId,
			},
			UseRemoteGateways: to.BoolPtr(false),
		},
	}
	err = CreateVirtualNetworkPeering(ctx, virtualNetworkName, virtualNetworkPeeringName, virtualNetworkPeeringParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network peering: % +v", err)
	}
	t.Logf("created virtual network peering")

	err = GetVirtualNetworkPeering(ctx, virtualNetworkName, virtualNetworkPeeringName)
	if err != nil {
		t.Fatalf("failed to get virtual network peering: %+v", err)
	}
	t.Logf("got virtual network peering")

	err = ListVirtualNetworkPeering(ctx, virtualNetworkName)
	if err != nil {
		t.Fatalf("failed to list virtual network peering: %+v", err)
	}
	t.Logf("listed virtual network peering")

	err = DeleteVirtualNetworkPeering(ctx, virtualNetworkName, virtualNetworkPeeringName)
	if err != nil {
		t.Fatalf("failed to delete virtual network peering: %+v", err)
	}
	t.Logf("deleted virtual network peering")

}
