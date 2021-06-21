// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	network "github.com/Azure-Samples/azure-sdk-for-go-samples/network/sdk"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVirtualMachineExtension(t *testing.T) {
	groupName := config.GenerateGroupName("compute")
	config.SetGroupName(groupName)

	virtualMachineExtensionName := config.AppendRandomSuffix("virtualmachineextension")
	networkInterfaceName := config.AppendRandomSuffix("networkinterface")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	subnetName := config.AppendRandomSuffix("subnet")
	virtualMachineName := config.AppendRandomSuffix("vm")
	ipConfigurationName := config.AppendRandomSuffix("ipconfiguration")

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
				AddressPrefixes: []*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = network.CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/16"),
		},
	}
	subNetID, err := network.CreateSubnet(ctx, virtualNetworkName, subnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	publicIpAddressId, err := network.CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	networkInterfaceParameters := armnetwork.NetworkInterface{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.NetworkInterfacePropertiesFormat{
			EnableAcceleratedNetworking: to.BoolPtr(true),
			IPConfigurations: []*armnetwork.NetworkInterfaceIPConfiguration{
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

	nicId, _, err := network.CreateNetworkInterface(ctx, networkInterfaceName, networkInterfaceParameters)
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
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
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
				DataDisks: []*armcompute.DataDisk{
					{
						CreateOption: armcompute.DiskCreateOptionTypesEmpty.ToPtr(),
						DiskSizeGB:   to.Int32Ptr(1023),
						Lun:          to.Int32Ptr(0),
					},
					{
						CreateOption: armcompute.DiskCreateOptionTypesEmpty.ToPtr(),
						DiskSizeGB:   to.Int32Ptr(1023),
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
						StorageAccountType: armcompute.StorageAccountTypesStandardLRS.ToPtr(),
					},
					Name: to.StringPtr("myVMosdisk"),
				},
			},
		},
	}

	_, err = CreateVirtualMachine(ctx, virtualMachineName, virtualMachineProbably)
	if err != nil {
		t.Fatalf("failed to create virtual machine: % +v", err)
	}

	extensionParameters := armcompute.VirtualMachineExtension{
		Resource: armcompute.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armcompute.VirtualMachineExtensionProperties{
			AutoUpgradeMinorVersion: to.BoolPtr(true),
			Publisher:               to.StringPtr("Microsoft.Azure.NetworkWatcher"),
			Type:                    to.StringPtr("NetworkWatcherAgentWindows"),
			TypeHandlerVersion:      to.StringPtr("1.4"),
		},
	}

	err = CreateVirtualMachineExtension(ctx, virtualMachineName, virtualMachineExtensionName, extensionParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine extension: % +v", err)
	}
	t.Logf("created virtual machine extension")

	err = GetVirtualMachineExtension(ctx, virtualMachineName, virtualMachineExtensionName)
	if err != nil {
		t.Fatalf("failed to get virtual machine extension: %+v", err)
	}
	t.Logf("got virtual machine extension")

	err = ListVirtualMachineExtension(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to list virtual machine extension: %+v", err)
	}
	t.Logf("listed virtual machine extension")

	virtualMachineExtensionUpdateParameters := armcompute.VirtualMachineExtensionUpdate{
		Properties: &armcompute.VirtualMachineExtensionUpdateProperties{
			AutoUpgradeMinorVersion: to.BoolPtr(true),
			Settings:                "{\"commandToExecute\": \"powershell.exe -c \"Get-Process | Where-Object { $_.CPU -lt 100 }\"}",
		},
	}
	err = UpdateVirtualMachineExtensionTags(ctx, virtualMachineName, virtualMachineExtensionName, virtualMachineExtensionUpdateParameters)
	if err != nil {
		t.Fatalf("failed to update virtual machine extension: %+v", err)
	}
	t.Logf("updated virtual machine extension tags")

	err = DeleteVirtualMachineExtension(ctx, virtualMachineName, virtualMachineExtensionName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine extension: %+v", err)
	}
	t.Logf("deleted virtual machine extension")

}
