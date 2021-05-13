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

func TestInterface(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	config.SetLocation("eastus")

	networkInterfaceName := config.AppendRandomSuffix("networkinterface")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	subnetName := config.AppendRandomSuffix("subnet")
	virtualMachineName := config.AppendRandomSuffix("vm")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.DefaultLocation())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateVirtualNetwork(ctx, virtualNetworkName)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}
	t.Logf("created virtual network")

	body := `{
		"addressPrefix": "10.0.0.0/16"
	  }
	`
	_, err = CreateSubnet(ctx, virtualNetworkName, subnetName, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	publicIPAddress := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	err = CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddress)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	nicId, err := CreateNetworkInterface(ctx, networkInterfaceName, publicIpAddressName, virtualNetworkName, subnetName)
	if err != nil {
		t.Fatalf("failed to create network interface: % +v", err)
	}
	t.Logf("created network interface")

	err = CreateVirtualMachine(ctx, virtualMachineName, nicId)
	if err != nil {
		t.Fatalf("failed to create virtual machine: % +v", err)
	}

	err = ListNetworkInterfaceIpConfiguration(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to list network interface ip configuration: %+v", err)
	}
	t.Logf("listed network interface ip configuration")

	err = ListNetworkInterfaceLoadBalancer(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to list network interface load balancer: %+v", err)
	}
	t.Logf("listed network interface load balancer")

	err = GetNetworkInterface(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to get network interface: %+v", err)
	}
	t.Logf("got network interface")

	err = ListNetworkInterface(ctx)
	if err != nil {
		t.Fatalf("failed to list network interface: %+v", err)
	}
	t.Logf("listed network interface")

	err = ListAllNetworkInterface(ctx)
	if err != nil {
		t.Fatalf("failed to list all network interface: %+v", err)
	}
	t.Logf("listed all network interface")

	err = BeginListEffectiveRouteTable(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to list all network security groups applied to a network interface: %+v", err)
	}
	t.Logf("listed all network security groups applied to a network interface")

	err = BeginGetEffectiveRouteTable(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to get all route tables applied to a network interface: %+v", err)
	}
	t.Logf("got all route tables applied to a network interface")

	err = UpdateNetworkInterfaceTags(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to update tags for network interface: %+v", err)
	}
	t.Logf("updated network interface tags")

	err = DeleteVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine: %+v", err)
	}

	err = DeleteNetworkInterface(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to delete network interface: %+v", err)
	}
	t.Logf("deleted network interface")

}
