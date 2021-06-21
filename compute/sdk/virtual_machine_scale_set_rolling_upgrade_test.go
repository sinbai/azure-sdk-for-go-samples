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

func TestVirtualMachineScaleSetRollingUpgrade(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	vmScaleSetName := config.AppendRandomSuffix("virtualmachinescaleset")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName := config.AppendRandomSuffix("subnet")

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
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
				RollingUpgradePolicy: &armcompute.RollingUpgradePolicy{
					MaxUnhealthyInstancePercent:         to.Int32Ptr(100),
					MaxUnhealthyUpgradedInstancePercent: to.Int32Ptr(100),
				},
			},
			Overprovision: to.BoolPtr(true),
			VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
				NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
					NetworkInterfaceConfigurations: []*armcompute.VirtualMachineScaleSetNetworkConfiguration{{
						Name: to.StringPtr("testPC"),
						Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
							EnableIPForwarding: to.BoolPtr(true),
							IPConfigurations: []*armcompute.VirtualMachineScaleSetIPConfiguration{{
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
						DiskSizeGB:   to.Int32Ptr(512),
						ManagedDisk: &armcompute.VirtualMachineScaleSetManagedDiskParameters{
							StorageAccountType: armcompute.StorageAccountTypesStandardLRS.ToPtr(),
						},
					},
				},
			},
		},

		SKU: &armcompute.SKU{
			Capacity: to.Int64Ptr(1),
			Name:     to.StringPtr("Standard_D1_v2"),
			Tier:     to.StringPtr("Standard"),
		},
	}
	err = CreateVirtualMachineScaleSet(ctx, vmScaleSetName, virtualMachineScaleSetParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine scale set: % +v", err)
	}

	err = StartRollingExtensionUpgrade(ctx, vmScaleSetName)
	if err != nil {
		t.Fatalf("failed to start virtual machine scale set rolling upgrade: %+v", err)
	}
	t.Logf("started virtual machine scale set rolling upgrade")

	err = StartRollingUpgradeOSUpgrade(ctx, vmScaleSetName, true)
	if err != nil {
		t.Fatalf("failed to start a rolling upgrade to move all virtual machine scale set instances to the latest available Platform Image OS version: %+v", err)
	}
	t.Logf("started a rolling upgrade to move all virtual machine scale set instances to the latest available Platform Image OS version")

	err = GetLatestVirtualMachineScaleSetRollingUpgrade(ctx, vmScaleSetName)
	if err != nil {
		t.Fatalf("failed to get the status of the latest virtual machine scale set rolling upgrade: %+v", err)
	}
	t.Logf("got the status of the latest virtual machine scale set rolling upgrade")

	err = StartRollingUpgradeOSUpgrade(ctx, vmScaleSetName, false)
	if err != nil {
		t.Fatalf("failed to start a rolling upgrade to move all virtual machine scale set instances to the latest available Platform Image OS version: %+v", err)
	}

	err = CancelScaleSetRollingUpgrade(ctx, vmScaleSetName)
	if err != nil {
		t.Fatalf("failed to cancel the current virtual machine scale set rolling upgrade: %+v", err)
	}
	t.Logf("cancelled the current virtual machine scale set rolling upgrade")
}
