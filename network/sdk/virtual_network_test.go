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

func TestVirtualNetwork(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
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
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}
	t.Logf("created virtual network")

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                     to.StringPtr("10.0.1.0/24"),
			PrivateLinkServiceNetworkPolicies: to.StringPtr("Disabled"),
		},
	}
	_, err = CreateSubnet(ctx, virtualNetworkName, subNetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	ipAddress := "10.0.1.4"
	err = CheckIPAddressAvailability(ctx, virtualNetworkName, ipAddress)
	if err != nil {
		t.Fatalf("failed to check ip address availability: %+v", err)
	}
	t.Logf("checked ip address availability")

	err = ListUsageVirtualNetwork(ctx, virtualNetworkName)
	if err != nil {
		t.Fatalf("failed to list usage virtual network: %+v", err)
	}
	t.Logf("listed usage virtual network")

	err = GetVirtualNetwork(ctx, virtualNetworkName)
	if err != nil {
		t.Fatalf("failed to get virtual network: %+v", err)
	}
	t.Logf("got virtual network")

	err = ListVirtualNetwork(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual network: %+v", err)
	}
	t.Logf("listed virtual network")

	err = ListAllVirtualNetwork(ctx)
	if err != nil {
		t.Fatalf("failed to list all virtual network: %+v", err)
	}
	t.Logf("listed all virtual network")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateVirtualNetworkTags(ctx, virtualNetworkName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for virtual network: %+v", err)
	}
	t.Logf("updated virtual network tags")

	err = DeleteVirtualNetwork(ctx, virtualNetworkName)
	if err != nil {
		t.Fatalf("failed to delete virtual network: %+v", err)
	}
	t.Logf("deleted virtual network")

}
