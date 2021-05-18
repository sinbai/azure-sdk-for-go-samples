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
)

func TestBastionHost(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	bastionHostName := config.AppendRandomSuffix("bastionhost")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subnetName := config.AppendRandomSuffix("subnet")

	bastionVirtualNetworkName := config.AppendRandomSuffix("bastionvirutalnetwork")
	bastionSubnetName := "AzureBastionSubnet"

	virtualMachineName := config.AppendRandomSuffix("virtualmachine")
	interfaceName := config.AppendRandomSuffix("interface")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateVirtualNetwork(ctx, virtualNetworkName, "10.0.0.0/16")
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	body := `{
		"addressPrefix": "10.0.0.0/24"
	  }
	`
	subNetId, err := CreateSubnet(ctx, virtualNetworkName, subnetName, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	nicId, err := createNetworkInterface(ctx, interfaceName, virtualNetworkName, subNetId)
	if err != nil {
		t.Fatalf("failed to create network interface: % +v", err)
	}

	err = CreateVirtualMachine(ctx, virtualMachineName, nicId)
	if err != nil {
		t.Fatalf("failed to create virtual machine: % +v", err)
	}

	publicIpAddressId, err := createPublicIPAddress(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	err = CreateVirtualNetwork(ctx, bastionVirtualNetworkName, "10.0.0.0/16")
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	body = `{
		"addressPrefix": "10.0.0.0/24"
	  }
	`
	bastionSubnetId, err := CreateSubnet(ctx, bastionVirtualNetworkName, bastionSubnetName, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	err = CreateBastionHost(ctx, bastionHostName, bastionSubnetId, publicIpAddressId)
	if err != nil {
		t.Fatalf("failed to create bastion host: % +v", err)
	}
	t.Logf("created bastion host")

	err = GetBastionHost(ctx, bastionHostName)
	if err != nil {
		t.Fatalf("failed to get bastion host: %+v", err)
	}
	t.Logf("got bastion host")

	err = ListBastionHost(ctx)
	if err != nil {
		t.Fatalf("failed to list bastion host: %+v", err)
	}
	t.Logf("listed bastion host")

	err = ListBastionHostByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to listbastion host by resource group: %+v", err)
	}
	t.Logf("listedbastion host by resource group")

	err = DeleteBastionHost(ctx, bastionHostName)
	if err != nil {
		t.Fatalf("failed to delete bastion host: %+v", err)
	}
	t.Logf("deleted bastion host")
}
