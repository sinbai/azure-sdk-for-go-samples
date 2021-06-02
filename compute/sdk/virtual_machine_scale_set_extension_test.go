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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVirtualMachineScaleSetExtension(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualMachineScaleSetExtensionName := config.AppendRandomSuffix("virtualmachinescalesetextension")
	virtualMachineScaleSetName := config.AppendRandomSuffix("virtualmachinescaleset")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName := config.AppendRandomSuffix("subnet")

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
	_, err = network.CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.1.0/24"),
		},
	}
	subnetId, err := network.CreateSubnet(ctx, virtualNetworkName, subNetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	virtualMachineScaleSetParameters := armcompute.VirtualMachineScaleSet{
		Resource: armcompute.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armcompute.VirtualMachineScaleSetProperties{
			UpgradePolicy: &armcompute.UpgradePolicy{
				Mode: armcompute.UpgradeModeManual.ToPtr(),
			},
			Overprovision: to.BoolPtr(true),
			VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
				NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
					NetworkInterfaceConfigurations: &[]armcompute.VirtualMachineScaleSetNetworkConfiguration{{
						Name: to.StringPtr("testPC"),
						Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
							EnableIPForwarding: to.BoolPtr(true),
							IPConfigurations: &[]armcompute.VirtualMachineScaleSetIPConfiguration{{
								Name: to.StringPtr("testPC"),
								Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
									Subnet: &armcompute.APIEntityReference{
										ID: &subnetId,
									},
								},
							}},
							Primary: to.BoolPtr(true),
						},
					}},
				},
				OSProfile: &armcompute.VirtualMachineScaleSetOSProfile{
					AdminPassword:      to.StringPtr("Aa!1()-xyz"),
					AdminUsername:      to.StringPtr("testuser"),
					ComputerNamePrefix: to.StringPtr("testPC"),
				},
				StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
					ImageReference: &armcompute.ImageReference{
						Offer:     to.StringPtr("WindowsServer"),
						Publisher: to.StringPtr("MicrosoftWindowsServer"),
						SKU:       to.StringPtr("2016-Datacenter"),
						Version:   to.StringPtr("latest"),
					},
					OSDisk: &armcompute.VirtualMachineScaleSetOSDisk{
						Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
						CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
						DiskSizeGb:   to.Int32Ptr(512),
						ManagedDisk: &armcompute.VirtualMachineScaleSetManagedDiskParameters{
							StorageAccountType: armcompute.StorageAccountTypesStandardLrs.ToPtr(),
						},
					},
				},
			},
		},

		SKU: &armcompute.SKU{
			Capacity: to.Int64Ptr(2),
			Name:     to.StringPtr("Standard_D1_v2"),
			Tier:     to.StringPtr("Standard"),
		},
	}
	err = CreateVirtualMachineScaleSet(ctx, virtualMachineScaleSetName, virtualMachineScaleSetParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine scale set: % +v", err)
	}

	extensionParameters := armcompute.VirtualMachineScaleSetExtension{
		Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
			AutoUpgradeMinorVersion: to.BoolPtr(true),
			Publisher:               to.StringPtr("Microsoft.Azure.NetworkWatcher"),
			Type:                    to.StringPtr("NetworkWatcherAgentWindows"),
			TypeHandlerVersion:      to.StringPtr("1.4"),
		},
	}
	err = CreateVirtualMachineScaleSetExtension(ctx, virtualMachineScaleSetName,
		virtualMachineScaleSetExtensionName, extensionParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine scale set extension: % +v", err)
	}
	t.Logf("created virtual machine scale set extension")

	err = GetVirtualMachineScaleSetExtension(ctx, virtualMachineScaleSetName, virtualMachineScaleSetExtensionName)
	if err != nil {
		t.Fatalf("failed to get virtual machine scale set extension: %+v", err)
	}
	t.Logf("got virtual machine scale set extension")

	err = ListVirtualMachineScaleSetExtension(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to list virtual machine scale set extension: %+v", err)
	}
	t.Logf("listed virtual machine scale set extension")

	virtualMachineScaleSetExtensionUpdateParameters := armcompute.VirtualMachineScaleSetExtensionUpdate{
		Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
			AutoUpgradeMinorVersion: to.BoolPtr(true),
		},
	}
	err = UpdateVirtualMachineScaleSetExtensionTags(ctx, virtualMachineScaleSetName,
		virtualMachineScaleSetExtensionName, virtualMachineScaleSetExtensionUpdateParameters)
	if err != nil {
		t.Fatalf("failed to update tags for virtual machine scale set extension: %+v", err)
	}
	t.Logf("updated virtual machine scale set extension tags")

	err = DeleteVirtualMachineScaleSetExtension(ctx, virtualMachineScaleSetName, virtualMachineScaleSetExtensionName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine scale set extension: %+v", err)
	}
	t.Logf("deleted virtual machine scale set extension")

}
