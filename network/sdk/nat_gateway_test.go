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

func TestNatGateWay(t *testing.T) {
	natGatewayName := config.AppendRandomSuffix("natgateway")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	prefixName := config.AppendRandomSuffix("pipprefix")

	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

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

	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	publicIPPrefixParameters := armnetwork.PublicIPPrefix{
		Resource: armnetwork.Resource{
			Name:     to.StringPtr(prefixName),
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
	publicIpPrefixId, err := CreatePublicIPPrefix(ctx, prefixName, publicIPPrefixParameters)
	if err != nil {
		t.Fatalf("failed to create public ip prefix: %+v", err)
	}

	natGatewayParameters := armnetwork.NatGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.NatGatewayPropertiesFormat{
			PublicIPAddresses: &[]*armnetwork.SubResource{
				{
					ID: &publicIpAddressId,
				},
			},
			PublicIPPrefixes: &[]*armnetwork.SubResource{
				{
					ID: &publicIpPrefixId,
				},
			},
		},
		SKU: &armnetwork.NatGatewaySKU{
			Name: armnetwork.NatGatewaySKUNameStandard.ToPtr(),
		},
	}
	err = CreateNatGateway(ctx, natGatewayName, natGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create nat gateway: %+v", err)
	}
	t.Logf("created nat gateway")

	err = GetNatGateway(ctx, natGatewayName)
	if err != nil {
		t.Fatalf("failed to get nat gateway: %+v", err)
	}
	t.Logf("got nat gateway")

	err = ListNatGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list nat gateway: %+v", err)
	}
	t.Logf("listed nat gateway")

	err = ListAllNatGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list all nat gateway: %+v", err)
	}
	t.Logf("listed all nat gateway")

	err = DeleteNatGatewayGroup(ctx, natGatewayName)
	if err != nil {
		t.Fatalf("failed to delete nat gateway: %+v", err)
	}
	t.Logf("deleted nat gateway")
}
