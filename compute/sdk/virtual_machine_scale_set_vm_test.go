// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	network "github.com/Azure-Samples/azure-sdk-for-go-samples/network/sdk"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVirtualMachineScaleSetVm(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualMachineScaleSetName := config.AppendRandomSuffix("virtualmachinescaleset")
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
				RollingUpgradePolicy: &armcompute.RollingUpgradePolicy{
					MaxUnhealthyInstancePercent:         to.Int32Ptr(100),
					MaxUnhealthyUpgradedInstancePercent: to.Int32Ptr(100),
				},
			},
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
			Capacity: to.Int64Ptr(1),
			Name:     to.StringPtr("Standard_D1_v2"),
			Tier:     to.StringPtr("Standard"),
		},
	}
	err = CreateVirtualMachineScaleSet(ctx, virtualMachineScaleSetName, virtualMachineScaleSetParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine scale set: % +v", err)
	}

	instanceId := 0
	for i := 0; i < 4; i++ {
		instanceId = i
		err = GetVirtualMachineScaleSetVmInstanceView(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
		if err != nil {
			if instanceId >= 3 {
				t.Fatalf("failed to redeploy a virtual machines in a VM scale set: %+v", err)
			}
			continue
		}
		break
	}
	t.Logf("got the status of a virtual machine from a VM scale set by instanceid: %+v", instanceId)

	// cannot use it successfully,error:"Operation 'performMaintenance' is not allowed since the Subscription of this VM is not eligible.
	err = VirtualMachineScaleSetVmPerformMaintenance(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to perform maintenance on a virtual machines in a VM scale set: % +v", err)
	}
	t.Logf("performed maintenance on a virtual machines in a VM scale set")

	err = RedeployVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to redeploy a virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("redeplied a virtual machines in a VM scale set")

	err = VirtualMachineScaleSetVmReimage(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to reimage a virtual machines in a VM scale set: % +v", err)
	}
	t.Logf("reimaged a virtual machines in a VM scale set")

	err = ReimageAllVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to reimage all the disks ( including data disks ) in a VM scale set instance: % +v", err)
	}
	t.Logf("reimaged all the disks ( including data disks ) in a VM scale set instance")

	err = ListVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to list all virtual machine scale set vm: %+v", err)
	}
	t.Logf("listed all virtual machine scale set vm")

	err = GetVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to get a virtual machine from a VM scale set: %+v", err)
	}
	t.Logf("got a virtual machine from a VM scale set")

	virtualMachineScaleSetVMParameters := armcompute.VirtualMachineScaleSetVM{
		Resource: armcompute.Resource{
			Tags: &map[string]string{"department": "HR"},
		},
	}
	err = UpdateVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId), virtualMachineScaleSetVMParameters)
	if err != nil {
		t.Fatalf("failed to update a virtual machine of a VM scale set: %+v", err)
	}
	t.Logf("updated a virtual machine of a VM scale set")

	err = RestartVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to restart a virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("restarted a virtual machines in a VM scale set")

	err = VirtualMachineScaleSetVmPowerOff(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to power off a virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("powered off a virtual machines in a VM scale set")

	err = StartVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to start a virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("start a virtual machines in a VM scale set")

	runCommandInputParameters := armcompute.RunCommandInput{
		CommandID: to.StringPtr("RunPowerShellScript"),
	}
	err = RunCommandOnVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId), runCommandInputParameters)
	if err != nil {
		t.Fatalf("failed to run command on a virtual machine in a VM scale set: %+v", err)
	}
	t.Logf("ran command on a virtual machine in a VM scale set")

	err = DeallocateVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to deallocte a virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("deallocted a virtual machines in a VM scale set")

	err = DeleteVirtualMachineScaleSetVm(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
	if err != nil {
		t.Fatalf("failed to delete a virtual machine from a VM scale set: %+v", err)
	}
	t.Logf("deleted a virtual machine from a VM scale set")
}
