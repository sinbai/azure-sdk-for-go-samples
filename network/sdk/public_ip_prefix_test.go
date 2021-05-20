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

func TestPublicIPPrefix(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	publicIpPrefixName := config.AppendRandomSuffix("pipperfix")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	publicIPPrefixPro := armnetwork.PublicIPPrefix{
		Resource: armnetwork.Resource{
			Name:     to.StringPtr(publicIpPrefixName),
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.PublicIPPrefixPropertiesFormat{
			PrefixLength:           to.Int32Ptr(30),
			PublicIPAddressVersion: armnetwork.IPVersionIPv4.ToPtr(),
		},
		SKU: &armnetwork.PublicIPPrefixSKU{
			Name: armnetwork.PublicIPPrefixSKUNameStandard.ToPtr(),
		},
	}
	_, err = CreatePublicIPPrefix(ctx, publicIpPrefixName, publicIPPrefixPro)
	if err != nil {
		t.Fatalf("failed to create public ip prefix: %+v", err)
	}
	t.Logf("created public ip prefix")

	err = GetPublicIPPrefix(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to get public ip prefix: %+v", err)
	}
	t.Logf("got public ip prefix")

	err = ListPublicIPPrefix(ctx)
	if err != nil {
		t.Fatalf("failed to list public ip prefix: %+v", err)
	}
	t.Logf("listed public ip prefix")

	err = ListAllPublicIPPrefix(ctx)
	if err != nil {
		t.Fatalf("failed to list all public ip prefix: %+v", err)
	}
	t.Logf("listed all public ip prefix")

	err = UpdatePublicIPPrefixTags(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to update tags for public ip prefix: %+v", err)
	}
	t.Logf("updated public ip prefix tags")

	err = DeletePublicIPPrefix(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to delete public ip prefix: %+v", err)
	}
	t.Logf("deleted public ip prefix")
}
