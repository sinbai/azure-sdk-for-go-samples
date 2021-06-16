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
	"github.com/Azure/go-autorest/autorest/to"
)

func TestNetworkVirtualAppliance(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	networkVirtualApplianceName := config.AppendRandomSuffix("networkvirtualappliance")

	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
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
	virtualWanId, err := CreateVirtualWan(ctx, virtualWanName, virtualWANParameters)
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
				ID: &virtualWanId,
			},
		},
	}

	virtualHubId, err := CreateVirtualHub(ctx, virtualHubName, virtualHubParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}

	parametersNetworkVirtualAppliance := armnetwork.NetworkVirtualAppliance{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.NetworkVirtualAppliancePropertiesFormat{
			NvaSKU: &armnetwork.VirtualApplianceSKUProperties{
				BundledScaleUnit:   to.StringPtr("3"),
				MarketPlaceVersion: to.StringPtr("17.4.0928"),
				Vendor:             to.StringPtr("ciscosdwan"),
			},
			VirtualApplianceAsn: to.Int64Ptr(10000),
			VirtualHub: &armnetwork.SubResource{
				ID: &virtualHubId,
			},
		},
	}

	err = CreateNetworkVirtualAppliance(ctx, networkVirtualApplianceName, parametersNetworkVirtualAppliance)
	if err != nil {
		t.Fatalf("failed to create network virtual appliance: % +v", err)
	}
	t.Logf("created network virtual appliance")

	err = GetNetworkVirtualAppliance(ctx, networkVirtualApplianceName)
	if err != nil {
		t.Fatalf("failed to get network virtual appliance: %+v", err)
	}
	t.Logf("got network virtual appliance")

	err = ListNetworkVirtualAppliance(ctx)
	if err != nil {
		t.Fatalf("failed to list network virtual appliance: %+v", err)
	}
	t.Logf("listed network virtual appliance")

	err = ListNetworkVirtualApplianceByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to lis tnetwork virtual appliance by resource group: %+v", err)
	}
	t.Logf("listed network virtual appliance by resource group")

	// Error Message: "The requested resource does not support http method 'PATCH'."
	// seems it’s API limitation, disable it for now.
	// tagsObjectParameters := armnetwork.TagsObject{
	// 	Tags: &map[string]*string{"key1": to.StringPtr("value1"), "key2": to.StringPtr("value2")},
	// }
	// err = UpdateNetworkVirtualApplianceTags(ctx, networkVirtualApplianceName, tagsObjectParameters)
	// if err != nil {
	// 	t.Fatalf("failed to update tags for network virtual appliance: %+v", err)
	// }
	// t.Logf("updated network virtual appliance tags")

	err = DeleteNetworkVirtualAppliance(ctx, networkVirtualApplianceName)
	if err != nil {
		t.Fatalf("failed to delete network virtual appliance: %+v", err)
	}
	t.Logf("deleted network virtual appliance")
}
