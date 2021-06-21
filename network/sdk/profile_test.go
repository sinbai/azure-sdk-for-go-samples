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
	storage "github.com/Azure-Samples/azure-sdk-for-go-samples/storage/sdk"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/marstr/randname"
)

func TestNetworkProfile(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	networkProfileName := config.AppendRandomSuffix("networkprofile")
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	storageAccountName := randname.Prefixed{Prefix: "storageaccount", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subnetName := config.AppendRandomSuffix("subnet")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWANParameters := armnetwork.VirtualWAN{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     map[string]*string{"key1": to.StringPtr("value1")},
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
			Tags:     map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualHubProperties{
			AddressPrefix: to.StringPtr("10.168.0.0/24"),
			SKU:           to.StringPtr("Basic"),
			VirtualWan: &armnetwork.SubResource{
				ID: &virtualWanID,
			},
		},
	}
	_, err = CreateVirtualHub(ctx, virtualHubName, virtualHubParameters, false)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}

	storageAccountCreateParameters := armstorage.StorageAccountCreateParameters{
		Kind:     armstorage.KindStorage.ToPtr(),
		Location: to.StringPtr(config.Location()),
		SKU: &armstorage.SKU{
			Name: armstorage.SKUNameStandardLRS.ToPtr(),
		},
	}
	_, err = storage.CreateStorageAccount(ctx, storageAccountName, storageAccountCreateParameters)
	if err != nil {
		t.Fatalf("failed to create storage account: % +v", err)
	}

	virtualNetworkParameters := armnetwork.VirtualNetwork{
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
			AddressPrefix: to.StringPtr("10.0.0.0/16"),
		},
	}
	subnetId1, err := CreateSubnet(ctx, virtualNetworkName, subnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create subnet: % +v", err)
	}

	networkProfileParameters := armnetwork.NetworkProfile{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.NetworkProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: []*armnetwork.ContainerNetworkInterfaceConfiguration{{
				Name: to.StringPtr("eth1"),
				Properties: &armnetwork.ContainerNetworkInterfaceConfigurationPropertiesFormat{
					IPConfigurations: []*armnetwork.IPConfigurationProfile{
						{
							Name: to.StringPtr("ipconfig1"),
							Properties: &armnetwork.IPConfigurationProfilePropertiesFormat{
								Subnet: &armnetwork.Subnet{
									SubResource: armnetwork.SubResource{
										ID: &subnetId1,
									},
								},
							},
						},
					},
				},
			}},
		},
	}

	err = CreateNetworkProfile(ctx, networkProfileName, networkProfileParameters)
	if err != nil {
		t.Fatalf("failed to create network profile: % +v", err)
	}
	t.Logf("created network profile")

	err = GetNetworkProfile(ctx, networkProfileName)
	if err != nil {
		t.Fatalf("failed to get network profile: %+v", err)
	}
	t.Logf("got network profile")

	err = ListNetworkProfile(ctx)
	if err != nil {
		t.Fatalf("failed to list network profile: %+v", err)
	}
	t.Logf("listed network profile")

	err = ListAllNetworkProfile(ctx)
	if err != nil {
		t.Fatalf("failed to list all network profile: %+v", err)
	}
	t.Logf("listed all network profile")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateNetworkProfileTags(ctx, networkProfileName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for network profile: %+v", err)
	}
	t.Logf("updated network profile tags")

	err = DeleteNetworkProfile(ctx, networkProfileName)
	if err != nil {
		t.Fatalf("failed to delete network profile: %+v", err)
	}
	t.Logf("deleted network profile")

}
