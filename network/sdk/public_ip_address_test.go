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

func TestPublicIPAddress(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	publicIpAddressName := config.AppendRandomSuffix("pipaddress")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Name:     to.StringPtr(publicIpAddressName),
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
		},
		SKU: &armnetwork.PublicIPAddressSKU{
			Name: armnetwork.PublicIPAddressSKUNameStandard.ToPtr(),
		},
	}

	_, err = CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}
	t.Logf("created public ip address")

	err = GetPublicIPAddress(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to get public ip address: %+v", err)
	}
	t.Logf("got public ip address")

	err = ListPublicIPAddress(ctx)
	if err != nil {
		t.Fatalf("failed to list public ip address: %+v", err)
	}
	t.Logf("listed public ip address")

	err = ListAllPublicIPAddress(ctx)
	if err != nil {
		t.Fatalf("failed to list all public ip address: %+v", err)
	}
	t.Logf("listed all public ip address")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdatePublicIPAddressTags(ctx, publicIpAddressName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for public ip address: %+v", err)
	}
	t.Logf("updated public ip address tags")

	err = DeletePublicIPAddress(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to delete public ip address: %+v", err)
	}
	t.Logf("deleted public ip address")
}
