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

func TestVirtualHubBgpConnection(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualHubBgpConnectionName := config.AppendRandomSuffix("virtualhubbgpconnection")
	ipConfigName := config.AppendRandomSuffix("ipconfiguration")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	subNetName := "RouteServerSubnet"
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualHubParameters := armnetwork.VirtualHub{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualHubProperties{
			SKU: to.StringPtr("Standard"),
		},
	}
	virtualHubId, err := CreateVirtualHub(ctx, virtualHubName, virtualHubParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}

	virtualNetworkParameters := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.5.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.5.1.0/24"),
		},
	}
	subnetId, err := CreateSubnet(ctx, virtualNetworkName, subNetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Name:     to.StringPtr(publicIpAddressName),
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
		},
	}

	publicIpId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	hubIPConfigurationParameters := armnetwork.HubIPConfiguration{
		SubResource: armnetwork.SubResource{
			ID: &virtualHubId,
		},
		Properties: &armnetwork.HubIPConfigurationPropertiesFormat{
			PrivateIPAddress:          to.StringPtr("10.5.1.18"),
			PrivateIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
			PublicIPAddress: &armnetwork.PublicIPAddress{
				Resource: armnetwork.Resource{
					ID: &publicIpId,
				},
			},
			Subnet: &armnetwork.Subnet{
				SubResource: armnetwork.SubResource{
					ID: &subnetId,
				},
			},
		},
	}
	err = CreateVirtualHubIp(ctx, virtualHubName, ipConfigName, hubIPConfigurationParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub ip: % +v", err)
	}

	bgpConnectionParameters := armnetwork.BgpConnection{
		Properties: &armnetwork.BgpConnectionProperties{
			PeerAsn: to.Int64Ptr(65514),
			PeerIP:  to.StringPtr("169.254.21.5"),
		},
	}
	err = CreateVirtualHubBgpConnection(ctx, virtualHubName, virtualHubBgpConnectionName, bgpConnectionParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub bgp connection: % +v", err)
	}
	t.Logf("created virtual hub bgp connection")

	err = GetVirtualHubBgpConnection(ctx, virtualHubName, virtualHubBgpConnectionName)
	if err != nil {
		t.Fatalf("failed to get virtual hub bgp connection: %+v", err)
	}
	t.Logf("got virtual hub bgp connection")

	err = DeleteVirtualHubBgpConnection(ctx, virtualHubName, virtualHubBgpConnectionName)
	if err != nil {
		t.Fatalf("failed to delete virtual hub bgp connection: %+v", err)
	}
	t.Logf("deleted virtual hub bgp connection")

}
