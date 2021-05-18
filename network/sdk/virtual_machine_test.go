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

func TestVirtualMachine(t *testing.T) {
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

	err = CreateVirtualNetwork(ctx, virtualNetworkName, "10.0.0.0/16")
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

	publicIPAddressPro := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressPro)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	nicId, err := CreateNetworkInterface(ctx, networkInterfaceName, publicIpAddressId, virtualNetworkName, subnetName)
	if err != nil {
		t.Fatalf("failed to create network interface: % +v", err)
	}
	t.Logf("created network interface")

	err = CreateVirtualMachine(ctx, virtualMachineName, nicId)
	if err != nil {
		t.Fatalf("failed to create virtual machine: % +v", err)
	}

	err = DeleteVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine: %+v", err)
	}
}
