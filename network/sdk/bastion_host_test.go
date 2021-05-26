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
	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/24"),
		},
	}
	subNetId, err := CreateSubnet(ctx, virtualNetworkName, subnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	networkInterfaceParameters := armnetwork.NetworkInterface{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.NetworkInterfacePropertiesFormat{
			IPConfigurations: &[]*armnetwork.NetworkInterfaceIPConfiguration{
				{
					Name: to.StringPtr("MyIpConfig"),
					Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.Subnet{SubResource: armnetwork.SubResource{ID: &subNetId}},
					},
				},
			},
		},
	}

	nicId, _, err := CreateNetworkInterface(ctx, interfaceName, networkInterfaceParameters)
	if err != nil {
		t.Fatalf("failed to create network interface: % +v", err)
	}

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
							ID: &nicId,
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

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			IdleTimeoutInMinutes:     to.Int32Ptr(4),
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

	virtualNetworkParameters = armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, bastionVirtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters = armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/24"),
		},
	}
	bastionSubnetId, err := CreateSubnet(ctx, bastionVirtualNetworkName, bastionSubnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	bastionHostParameters := armnetwork.BastionHost{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.BastionHostPropertiesFormat{
			IPConfigurations: &[]*armnetwork.BastionHostIPConfiguration{{
				Name: to.StringPtr("bastionHostIpConfiguration"),
				Properties: &armnetwork.BastionHostIPConfigurationPropertiesFormat{
					PublicIPAddress: &armnetwork.SubResource{
						ID: &publicIpAddressId,
					},
					Subnet: &armnetwork.SubResource{
						ID: &bastionSubnetId,
					},
				},
			}},
		},
	}
	err = CreateBastionHost(ctx, bastionHostName, bastionHostParameters)
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
