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

func TestVirtualNetworkGateway(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualNetworkGatewayName := config.AppendRandomSuffix("virtualnetworkgateway")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	gatewaySubNetName := "GatewaySubnet"
	ipConfigName := config.AppendRandomSuffix("ipconfig")

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	publicIPAddressPro := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Name:     to.StringPtr(publicIpAddressName),
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
			IdleTimeoutInMinutes:     to.Int32Ptr(4),
		},
	}

	err = CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressPro)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	err = CreateVirtualNetwork(ctx, virtualNetworkName, "10.0.0.0/16")
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	body := `{
		"addressPrefix": "10.0.1.0/24"
		}`
	_, err = CreateSubnet(ctx, virtualNetworkName, gatewaySubNetName, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	err = CreateVirtualNetworkGateway(ctx, virtualNetworkName, virtualNetworkGatewayName, publicIpAddressName, ipConfigName, gatewaySubNetName)
	if err != nil {
		t.Fatalf("failed to create virtual network gateway: % +v", err)
	}
	t.Logf("created virtual network gateway")

	err = ListVirtualNetworkGatewayConnections(ctx, virtualNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to list virtual network gateway connection: %+v", err)
	}
	t.Logf("got virtual network gateway connection")

	err = GetVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to get virtual network gateway: %+v", err)
	}
	t.Logf("got virtual network gateway")

	err = ListVirtualNetworkGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual network gateway: %+v", err)
	}
	t.Logf("listed virtual network gateway")

	peer := "10.0.0.2"
	err = BeginGetVirtualNetworkGatewayAdvertisedRoute(ctx, virtualNetworkGatewayName, peer)
	if err != nil {
		t.Fatalf("failed to begin get virtual network gateway advertised route: %+v", err)
	}
	t.Logf("began get virtual network gateway advertised route")

	err = BeginGetVirtualNetworkGatewayBgpPeerStatus(ctx, virtualNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to begin get bgp peer status: %+v", err)
	}
	t.Logf("began get virtual network gateway bgp peer status")

	err = BeginGetVirtualNetworkGatewayLearnedRoutes(ctx, virtualNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to begin get learned route: %+v", err)
	}
	t.Logf("began get virtual network gateway learned route")

	err = BeginVirtualNetworkGatewayReset(ctx, virtualNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to begin reset: %+v", err)
	}
	t.Logf("began reset virtual network gateway")

	err = UpdateVirtualNetworkGatewayTags(ctx, virtualNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to update tags for virtual network gateway: %+v", err)
	}
	t.Logf("updated virtual network gateway tags")

	for i := 0; i < 4; i++ {
		err = DeleteVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
		if err != nil {
			//t.Fatalf("failed to delete virtual network gateway: %+v", err)
			err = DeleteVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
			if err != nil {
				err = DeleteVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
				if err != nil {
					err = DeleteVirtualNetworkGateway(ctx, virtualNetworkGatewayName)
					if err != nil {
						t.Fatalf("failed to delete virtual network gateway: %+v", err)
					}
				}
			}
		}
	}

	t.Logf("deleted virtual network gateway")

}
