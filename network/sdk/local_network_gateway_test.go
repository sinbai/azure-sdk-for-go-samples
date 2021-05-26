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

func TestLocalNetworkGateway(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	localNetworkGatewayName := config.AppendRandomSuffix("localnetworkgateway")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	localNetworkGatewayParameters := armnetwork.LocalNetworkGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
			GatewayIPAddress: to.StringPtr("11.12.13.14"),
			LocalNetworkAddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.1.0.0/16")},
			},
		},
	}
	_, err = CreateLocalNetworkGateway(ctx, localNetworkGatewayName, localNetworkGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create local network gateway: % +v", err)
	}
	t.Logf("created local network gateway")

	err = GetLocalNetworkGateway(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to get local network gateway: %+v", err)
	}
	t.Logf("got local network gateway")

	err = ListLocalNetworkGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list local network gateway: %+v", err)
	}
	t.Logf("listed local network gateway")

	err = UpdateLocalNetworkGatewayTags(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to update tags for local network gateway: %+v", err)
	}
	t.Logf("updated local network gateway tags")

	err = DeleteLocalNetworkGateway(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to delete local network gateway: %+v", err)
	}
	t.Logf("deleted local network gateway")

}
