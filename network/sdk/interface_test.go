// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"

	compute "github.com/Azure-Samples/azure-sdk-for-go-samples/compute/sdk"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
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
	ipConfigurationName := config.AppendRandomSuffix("ipconfiguration")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.DefaultLocation())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualNetworkParameters := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}
	t.Logf("created virtual network")

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/16"),
		},
	}
	subNetID, err := CreateSubnet(ctx, virtualNetworkName, subnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	networkInterfaceParameters := armnetwork.NetworkInterface{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.NetworkInterfacePropertiesFormat{
			EnableAcceleratedNetworking: to.BoolPtr(true),
			IPConfigurations: &[]*armnetwork.NetworkInterfaceIPConfiguration{
				{
					Name: &ipConfigurationName,
					Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							Resource: armnetwork.Resource{
								ID: &publicIpAddressId,
							},
						},
						Subnet: &armnetwork.Subnet{SubResource: armnetwork.SubResource{ID: &subNetID}},
					},
				},
			},
		},
	}

	interfaceId, _, err := CreateNetworkInterface(ctx, networkInterfaceName, networkInterfaceParameters)
	if err != nil {
		t.Fatalf("failed to create network interface: % +v", err)
	}
	t.Logf("created network interface")

	virtualMachineProbably := armcompute.VirtualMachine{
		Resource: armcompute.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armcompute.VirtualMachineProperties{
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: armcompute.VirtualMachineSizeTypesStandardD2V2.ToPtr(),
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: &[]armcompute.NetworkInterfaceReference{
					{
						SubResource: armcompute.SubResource{
							ID: &interfaceId,
						},
						Properties: &armcompute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
			OSProfile: &armcompute.OSProfile{
				AdminPassword: to.StringPtr("Aa1!zyx_"),
				AdminUsername: to.StringPtr("testuser"),
				ComputerName:  to.StringPtr("myVM"),
				WindowsConfiguration: &armcompute.WindowsConfiguration{
					EnableAutomaticUpdates: to.BoolPtr(true),
				},
			},
			StorageProfile: &armcompute.StorageProfile{
				DataDisks: &[]armcompute.DataDisk{
					{
						CreateOption: armcompute.DiskCreateOptionTypesEmpty.ToPtr(),
						DiskSizeGb:   to.Int32Ptr(1023),
						Lun:          to.Int32Ptr(0),
					},
					{
						CreateOption: armcompute.DiskCreateOptionTypesEmpty.ToPtr(),
						DiskSizeGb:   to.Int32Ptr(1023),
						Lun:          to.Int32Ptr(1),
					},
				},
				ImageReference: &armcompute.ImageReference{
					Offer:     to.StringPtr("WindowsServer"),
					Publisher: to.StringPtr("MicrosoftWindowsServer"),
					SKU:       to.StringPtr("2016-Datacenter"),
					Version:   to.StringPtr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
					CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: armcompute.StorageAccountTypesStandardLrs.ToPtr(),
					},
					Name: to.StringPtr("myVMosdisk"),
				},
			},
		},
	}

	_, err = compute.CreateVirtualMachine(ctx, virtualMachineName, virtualMachineProbably)
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

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateNetworkInterfaceTags(ctx, networkInterfaceName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for network interface: %+v", err)
	}
	t.Logf("updated network interface tags")

	err = compute.DeleteVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine: %+v", err)
	}

	err = DeleteNetworkInterface(ctx, networkInterfaceName)
	if err != nil {
		t.Fatalf("failed to delete network interface: %+v", err)
	}
	t.Logf("deleted network interface")

}
