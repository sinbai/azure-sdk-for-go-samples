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

func TestVirtualNetworkGatewayConnection(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualNetworkGatewayConnectionName := config.AppendRandomSuffix("virtualnetworkgatewayconnection")
	localNetworkGatewayName := config.AppendRandomSuffix("localvirtualnetworkgateway")
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

	_, err = CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressPro)
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

	err = CreateLocalNetworkGateway(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to create local network gateway: % +v", err)
	}

	err = CreateVirtualNetworkGatewayConnection(ctx, virtualNetworkName, virtualNetworkGatewayConnectionName, virtualNetworkGatewayName,
		localNetworkGatewayName, publicIpAddressName, gatewaySubNetName, ipConfigName)
	if err != nil {
		t.Fatalf("failed to create virtual network gateway connection: % +v", err)
	}
	t.Logf("created virtual network gateway connection")

	err = BeginSetVirtualNetworkGatewayConnectionSharedKey(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to set the virtual network gateway connection shared key: %+v", err)
	}
	t.Logf("set the virtual network gateway connection shared key")

	err = GetVirtualNetworkGatewayConnectionSharedKey(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to get the virtual network gateway connection shared key: %+v", err)
	}
	t.Logf("got the virtual network gateway connection shared key")

	err = GetVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to get virtual network gateway connection: %+v", err)
	}
	t.Logf("got virtual network gateway connection")

	err = ListVirtualNetworkGatewayConnection(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual network gateway connection: %+v", err)
	}
	t.Logf("listed virtual network gateway connection")

	//need to sleep for a period of time to run successfully, otherwise "Another operation on this or dependent resource is in progress." will be reported.
	time.Sleep(time.Duration(60) * time.Second)

	err = BeginResetVirtualNetworkGatewayConnectionSharedKey(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to begin reset the virtual network gateway connection shared key: %+v", err)
	}
	t.Logf("began reset the virtual network gateway connection shared key")

	//need to sleep for a period of time to run successfully, otherwise "Another operation on this or dependent resource is in progress." will be reported.
	time.Sleep(time.Duration(60) * time.Second)

	err = UpdateVirtualNetworkGatewayConnectionTags(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to update tags for virtual network gateway connection: %+v", err)
	}
	t.Logf("updated virtual network gateway connection tags")

	err = DeleteVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to delete virtual network gateway connection: %+v", err)
	}
	t.Logf("deleted virtual network gateway connection")
}
