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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2021-01-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/marstr/randname"
)

func TestPacketCapture(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	packetCaptureName := config.AppendRandomSuffix("packetcapture")
	networkWatcherName := config.AppendRandomSuffix("networkwatcher")
	storageAccountName := randname.Prefixed{Prefix: "storageaccount", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()
	virtualMachineName := config.AppendRandomSuffix("virtualmachine")
	ipConfigurationName := config.AppendRandomSuffix("ipconfiguration")
	networkInterfaceName := config.AppendRandomSuffix("interface")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subnetName := config.AppendRandomSuffix("subnet")
	virtualMachineExtensionName := config.AppendRandomSuffix("virtualmachineextension")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateNetworkWatcher(ctx, networkWatcherName)
	if err != nil {
		t.Fatalf("failed to create network watcher: % +v", err)
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
	subnetId, err := CreateSubnet(ctx, virtualNetworkName, subnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create subnet: % +v", err)
	}

	networkInterfaceParameters := armnetwork.NetworkInterface{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.NetworkInterfacePropertiesFormat{
			IPConfigurations: &[]*armnetwork.NetworkInterfaceIPConfiguration{
				{
					Name: &ipConfigurationName,
					Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.Subnet{SubResource: armnetwork.SubResource{ID: &subnetId}},
					},
				},
			},
		},
	}

	nicId, _, err := CreateNetworkInterface(ctx, networkInterfaceName, networkInterfaceParameters)
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

	vmId, err := compute.CreateVirtualMachine(ctx, virtualMachineName, virtualMachineProbably)
	if err != nil {
		t.Fatalf("failed to create virtual machine: % +v", err)
	}

	storageAccountCreateParameters := armstorage.StorageAccountCreateParameters{
		Kind:     armstorage.KindStorage.ToPtr(),
		Location: to.StringPtr(config.Location()),
		SKU: &armstorage.SKU{
			Name: armstorage.SKUNameStandardLRS.ToPtr(),
		},
	}
	stroageAountId, err := CreateStorageAccount(ctx, storageAccountName, storageAccountCreateParameters)
	if err != nil {
		t.Fatalf("failed to create storage account: % +v", err)
	}

	virtualMachineExtensionParameters := armcompute.VirtualMachineExtension{
		Resource: armcompute.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armcompute.VirtualMachineExtensionProperties{
			AutoUpgradeMinorVersion: to.BoolPtr(true),
			Publisher:               to.StringPtr("Microsoft.Azure.NetworkWatcher"),
			Type:                    to.StringPtr("NetworkWatcherAgentWindows"),
			TypeHandlerVersion:      to.StringPtr("1.4"),
		},
	}
	err = compute.CreateVirtualMachineExtension(ctx, virtualMachineName, virtualMachineExtensionName, virtualMachineExtensionParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine extension: % +v", err)
	}

	packetCaptureParameters := armnetwork.PacketCapture{
		Properties: &armnetwork.PacketCaptureParameters{
			StorageLocation: &armnetwork.PacketCaptureStorageLocation{
				StorageID:   &stroageAountId,
				StoragePath: to.StringPtr("https://" + storageAccountName + ".blob.core.windows.net/capture/pc1.cap"),
			},
			Target: &vmId,
		},
	}
	err = CreatePacketCapture(ctx, networkWatcherName, packetCaptureName, packetCaptureParameters)
	if err != nil {
		t.Fatalf("failed to create packet capture: % +v", err)
	}
	t.Logf("created packet capture")

	err = GetPacketCapture(ctx, networkWatcherName, packetCaptureName)
	if err != nil {
		t.Fatalf("failed to get packet capture: %+v", err)
	}
	t.Logf("got packet capture")

	err = ListPacketCapture(ctx, networkWatcherName)
	if err != nil {
		t.Fatalf("failed to list packet capture: %+v", err)
	}
	t.Logf("listed packet capture")

	err = GetPacketCaptureStatus(ctx, networkWatcherName, packetCaptureName)
	if err != nil {
		t.Fatalf("failed to get packet capture status session: %+v", err)
	}
	t.Logf("got packet capture status session")

	err = StopPacketCapture(ctx, networkWatcherName, packetCaptureName)
	if err != nil {
		t.Fatalf("failed to stop packet capture: %+v", err)
	}
	t.Logf("stopped packet capture session")

	err = DeletePacketCapture(ctx, networkWatcherName, packetCaptureName)
	if err != nil {
		t.Fatalf("failed to delete packet capture: %+v", err)
	}
	t.Logf("deleted packet capture")

}
