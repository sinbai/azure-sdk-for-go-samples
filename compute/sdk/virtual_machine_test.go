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
	storage "github.com/Azure-Samples/azure-sdk-for-go-samples/storage/sdk"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2021-01-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/marstr/randname"
)

func TestVirtualMachine(t *testing.T) {
	groupName := config.GenerateGroupName("compute")
	config.SetGroupName(groupName)

	networkInterfaceName := config.AppendRandomSuffix("networkinterface")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	subnetName := config.AppendRandomSuffix("subnet")
	virtualMachineName := config.AppendRandomSuffix("vm")
	ipConfigurationName := config.AppendRandomSuffix("ipconfiguration")
	storageAccountName := randname.Prefixed{Prefix: "storageaccount", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()
	containerName := randname.Prefixed{Prefix: "blobcontainer", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()
	diskName := config.AppendRandomSuffix("disk")

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

	nicId, _, err := network.CreateNetworkInterface(ctx, networkInterfaceName, networkInterfaceParameters)
	if err != nil {
		t.Fatalf("failed to create network interface: % +v", err)
	}

	storageAccountCreateParameters := armstorage.StorageAccountCreateParameters{
		Kind:     armstorage.KindStorage.ToPtr(),
		Location: to.StringPtr(config.Location()),
		SKU: &armstorage.SKU{
			Name: armstorage.SKUNameStandardLRS.ToPtr(),
		},
	}
	_, err = storage.CreateStorageAccount(ctx, storageAccountName, storageAccountCreateParameters)
	if err != nil {
		t.Fatalf("failed to create storage account: % +v", err)
	}

	blobContainerParameters := armstorage.BlobContainer{}
	_, err = storage.CreateBlobContainer(ctx, storageAccountName, containerName, blobContainerParameters)
	if err != nil {
		t.Fatalf("failed to create blob container: % +v", err)
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
				ImageReference: &armcompute.ImageReference{
					Offer:     to.StringPtr("WindowsServer"),
					Publisher: to.StringPtr("MicrosoftWindowsServer"),
					SKU:       to.StringPtr("2016-Datacenter"),
					Version:   to.StringPtr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					OSType:       armcompute.OperatingSystemTypesWindows.ToPtr(),
					Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
					CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
					Name:         to.StringPtr("myVMosdisk"),
					Vhd: &armcompute.VirtualHardDisk{
						URI: to.StringPtr("http://" + storageAccountName + ".blob.core.windows.net/" + containerName + "/" + diskName + ".vhd"),
					},
				},
			},
			EvictionPolicy: armcompute.VirtualMachineEvictionPolicyTypesDeallocate.ToPtr(),
			BillingProfile: &armcompute.BillingProfile{
				MaxPrice: to.Float64Ptr(1),
			},
			Priority: armcompute.VirtualMachinePriorityTypesSpot.ToPtr(),
		},
	}

	_, err = CreateVirtualMachine(ctx, virtualMachineName, virtualMachineProbably)
	if err != nil {
		t.Fatalf("failed to create virtual machine: % +v", err)
	}
	t.Logf("created virtual machine")

	// Do not test from feedback
	// cannot use it successfully,error:"Operation 'performMaintenance' is not allowed since the Subscription of this VM is not eligible.
	// err = PerformMaintenanceVirtualMachine(ctx, virtualMachineName)
	// if err != nil {
	// 	t.Fatalf("failed to perform maintenance on a virtual machine: %+v", err)
	// }
	// t.Logf("performed maintenance on a virtual machine")

	// After synced with service team for ReimageVirtualMachine, currently they don’t support single VM. So we don’t need to add test case against it now.
	// cannot use it successfully, error: "The Reimage and OSUpgrade Virtual Machine actions require that the virtual machine has Automatic OS Upgrades enabled.
	// err = ReimageVirtualMachine(ctx, virtualMachineName)
	// if err != nil {
	// 	t.Fatalf("failed to reimage the virtual machine: %+v", err)
	// }
	// t.Logf("reimaged the virtual machine")

	err = InstanceVirtualMachineView(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to retrieve information about the run-time state of a virtual machine: %+v", err)
	}
	t.Logf("retrieved information about the run-time state of a virtual machine")

	err = ListVirtualMachineAvailableSizes(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to list all available virtual machine sizes: %+v", err)
	}
	t.Logf("listed all available virtual machine sizes")

	err = GetVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to get virtual machine: %+v", err)
	}
	t.Logf("got virtual machine")

	err = ListVirtualMachine(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual machine: %+v", err)
	}
	t.Logf("listed virtual machine")

	err = ListAllVirtualMachine(ctx)
	if err != nil {
		t.Fatalf("failed to list all virtual machine: %+v", err)
	}
	t.Logf("listed all virtual machine")

	err = ListVirtualMachineByLocation(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual machine by location: %+v", err)
	}
	t.Logf("listed virtual machine by location")

	runCommandInputParameters := armcompute.RunCommandInput{
		CommandID: to.StringPtr("RunPowerShellScript"),
	}
	err = RunCommandOnVirtualMachine(ctx, virtualMachineName, runCommandInputParameters)
	if err != nil {
		t.Fatalf("failed to run command on vm: %+v", err)
	}
	t.Logf("ran command on vm")

	err = RestartVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to restart virtual machine: %+v", err)
	}
	t.Logf("restarted virtual machine")

	err = VirtualMachinePowerOff(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to power off a virtual machine: %+v", err)
	}
	t.Logf("powered off virtual machine")

	err = StartVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to start virtual machine: %+v", err)
	}
	t.Logf("started virtual machine")

	err = VirtualMachinePowerOff(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to power off a virtual machine: %+v", err)
	}

	err = ReapplyVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to reapply virtual machine: %+v", err)
	}
	t.Logf("reapplied virtual machine")

	err = RedeployVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to redepoly virtual machine: %+v", err)
	}
	t.Logf("redepolied virtual machine")

	virtualMachineUpdateParameters := armcompute.VirtualMachineUpdate{
		Properties: &armcompute.VirtualMachineProperties{
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: &[]armcompute.NetworkInterfaceReference{{
					SubResource: armcompute.SubResource{
						ID: &nicId,
					},
					Properties: &armcompute.NetworkInterfaceReferenceProperties{
						Primary: to.BoolPtr(true),
					},
				}},
			},
		},
	}
	err = UpdateVirtualMachineTags(ctx, virtualMachineName, virtualMachineUpdateParameters)
	if err != nil {
		t.Fatalf("failed to update tags for virtual machine: %+v", err)
	}
	t.Logf("updated virtual machine tags")

	err = DeallocateVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to deallocate virtual machine: %+v", err)
	}
	t.Logf("deallocated virtual machine")

	err = ConvertVirtualMachineToManagedDisk(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to convert virtual machine disks from blob-based to managed disks: %+v", err)
	}
	t.Logf("converted virtual machine disks from blob-based to managed disks")

	err = SimulateEvictionVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to simulate the eviction of spot virtual machine: %+v", err)
	}
	t.Logf("simulated the eviction of spot virtual machine")

	err = DeallocateVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to deallocate virtual machine: %+v", err)
	}
	t.Logf("deallocated virtual machine")

	err = GenerializeVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to generialize virtual machine: %+v", err)
	}
	t.Logf("generialized virtual machine")

	err = DeleteVirtualMachine(ctx, virtualMachineName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine: %+v", err)
	}
	t.Logf("deleted virtual machine")
}
