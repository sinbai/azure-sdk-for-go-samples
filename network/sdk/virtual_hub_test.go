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

func TestVirtualHub(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")

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

	_, err = CreateVirtualHub(ctx, virtualHubName, virtualHubParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}
	t.Logf("created virtual hub")

	err = GetVirtualHub(ctx, virtualHubName)
	if err != nil {
		t.Fatalf("failed to get virtual hub: %+v", err)
	}
	t.Logf("got virtual hub")

	err = ListVirtualHub(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual hub: %+v", err)
	}
	t.Logf("listed virtual hub")

	err = ListVirtualHubByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual hub by resource group: %+v", err)
	}
	t.Logf("listed virtual hub by resource group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"key1": to.StringPtr("value1"), "key2": to.StringPtr("value2")},
	}
	err = UpdateVirtualHubTags(ctx, virtualHubName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for virtual hub: %+v", err)
	}
	t.Logf("updated virtual hub tags")

	// Error Message: The specified operation 'DeleteVirtualHub' is not supported.
	// Deletion is not supported when RoutingStatus on Hub is 'Provisioning'. Retry when state is not Provisioning
	// waiting for routing status to be provisioned and deleting
	time.Sleep(700 * time.Second)

	err = DeleteVirtualHub(ctx, virtualHubName)
	if err != nil {
		t.Fatalf("failed to delete virtual hub: %+v", err)
	}
	t.Logf("deleted virtual hub")

}
