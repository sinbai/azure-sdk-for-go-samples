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
	"github.com/Azure/go-autorest/autorest/to"
)

func TestConnectionMonitor(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	config.SetLocation("eastus")

	connectionMonitorName := config.AppendRandomSuffix("connectionmonitor")
	networkWatcherName := config.AppendRandomSuffix("networkwatcher")
	virtualMachineName := config.AppendRandomSuffix("virtualmachine")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subnetName := config.AppendRandomSuffix("subnet")
	ipConfigurationName := config.AppendRandomSuffix("ipconfiguration")
	networkInterfaceName := config.AppendRandomSuffix("interface")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.Location())

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

	connectionMonitorParameters := armnetwork.ConnectionMonitor{
		Location: to.StringPtr(config.Location()),
		Properties: &armnetwork.ConnectionMonitorParameters{
			Endpoints: &[]*armnetwork.ConnectionMonitorEndpoint{
				{
					Name:       to.StringPtr("vm1"),
					ResourceID: &vmId,
				},
				{
					Address: to.StringPtr("bing.com"),
					Name:    to.StringPtr("bing"),
				},
				{
					Address: to.StringPtr("google.com"),
					Name:    to.StringPtr("google"),
				},
			},
			TestConfigurations: &[]*armnetwork.ConnectionMonitorTestConfiguration{{
				Name:     to.StringPtr("testConfig1"),
				Protocol: armnetwork.ConnectionMonitorTestConfigurationProtocolTCP.ToPtr(),
				TCPConfiguration: &armnetwork.ConnectionMonitorTCPConfiguration{
					DisableTraceRoute: to.BoolPtr(true),
					Port:              to.Int32Ptr(80),
				},
				TestFrequencySec: to.Int32Ptr(60),
			}},
			TestGroups: &[]*armnetwork.ConnectionMonitorTestGroup{{
				Destinations:       &[]*string{to.StringPtr("bing"), to.StringPtr("google")},
				Disable:            to.BoolPtr(true),
				Name:               to.StringPtr("test1"),
				Sources:            &[]*string{to.StringPtr("vm1")},
				TestConfigurations: &[]*string{to.StringPtr("testConfig1")},
			}},
		},
	}
	err = CreateConnectionMonitor(ctx, networkWatcherName, connectionMonitorName, connectionMonitorParameters)
	if err != nil {
		t.Fatalf("failed to create connection monitor: % +v", err)
	}
	t.Logf("created connection monitor")

	err = GetConnectionMonitor(ctx, networkWatcherName, connectionMonitorName)
	if err != nil {
		t.Fatalf("failed to get connection monitor: %+v", err)
	}
	t.Logf("got connection monitor")

	err = ListConnectionMonitor(ctx, networkWatcherName)
	if err != nil {
		t.Fatalf("failed to list connection monitor: %+v", err)
	}
	t.Logf("listed connection monitor")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateConnectionMonitorTags(ctx, networkWatcherName, connectionMonitorName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for connection monitor: %+v", err)
	}
	t.Logf("updated connection monitor tags")

	err = DeleteConnectionMonitor(ctx, networkWatcherName, connectionMonitorName)
	if err != nil {
		t.Fatalf("failed to delete connection monitor: %+v", err)
	}
	t.Logf("deleted connection monitor")
}
