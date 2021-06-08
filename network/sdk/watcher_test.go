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

func TestNetworkWatcher(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	networkWatcherName := config.AppendRandomSuffix("networkwatcher")
	virtualMachineName := config.AppendRandomSuffix("virtualmachine")
	virtualNetworkNameGateway := config.AppendRandomSuffix("virtualnetwork")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	virtualNetworkGatewayName := config.AppendRandomSuffix("virtualnetworkgateway")
	publicIpAddressName := config.AppendRandomSuffix("publicipaddress")
	subnetGatewayName := "GatewaySubnet"
	subnetName := config.AppendRandomSuffix("subnet")
	storageAccountName := randname.Prefixed{Prefix: "storageaccount", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()
	ipConfigurationName := config.AppendRandomSuffix("ipconfiguration")
	networkInterfaceName := config.AppendRandomSuffix("interface")
	networkSecurityGroupName := config.AppendRandomSuffix("networksecuritygroup")
	virtualMachineExtensionName := config.AppendRandomSuffix("virtualmachineextension")

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
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
	t.Logf("created network watcher")

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
	_, err = CreateVirtualNetwork(ctx, virtualNetworkNameGateway, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/24"),
		},
	}
	subnetGatewayId, err := CreateSubnet(ctx, virtualNetworkNameGateway, subnetGatewayName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create subnet gateway: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			IdleTimeoutInMinutes:     to.Int32Ptr(10),
			PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
		},
	}
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	virtualNetWorkGatewayParameters := armnetwork.VirtualNetworkGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
			Active: to.BoolPtr(false),
			BgpSettings: &armnetwork.BgpSettings{
				Asn:               to.Int64Ptr(65515),
				BgpPeeringAddress: to.StringPtr("10.0.1.30"),
				PeerWeight:        to.Int32Ptr(0),
			},
			CustomRoutes: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("101.168.0.6/32")},
			},
			EnableBgp:           to.BoolPtr(false),
			EnableDNSForwarding: to.BoolPtr(false),
			GatewayType:         armnetwork.VirtualNetworkGatewayTypeVPN.ToPtr(),
			IPConfigurations: &[]*armnetwork.VirtualNetworkGatewayIPConfiguration{{
				Name: &ipConfigurationName,
				Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
					PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
					PublicIPAddress: &armnetwork.SubResource{
						ID: &publicIpAddressId,
					},
					Subnet: &armnetwork.SubResource{
						ID: &subnetGatewayId,
					},
				},
			}},
			SKU: &armnetwork.VirtualNetworkGatewaySKU{
				Name: armnetwork.VirtualNetworkGatewaySKUNameVPNGw1.ToPtr(),
				Tier: armnetwork.VirtualNetworkGatewaySKUTierVPNGw1.ToPtr(),
			},
			VPNType: armnetwork.VPNTypeRouteBased.ToPtr(),
		},
	}
	gatewayId, err := CreateVirtualNetworkGateway(ctx, virtualNetworkGatewayName, virtualNetWorkGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network gateway: % +v", err)
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
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters = armnetwork.Subnet{
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

	nicId, ipConfigPro, err := CreateNetworkInterface(ctx, networkInterfaceName, networkInterfaceParameters)
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

	networkSecurityGroupParameters := armnetwork.NetworkSecurityGroup{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	securityGroupId, err := CreateNetworkSecurityGroup(ctx, networkSecurityGroupName, networkSecurityGroupParameters)
	if err != nil {
		t.Fatalf("failed to create network security group: %+v", err)
	}

	troubleshootingParameters := armnetwork.TroubleshootingParameters{
		Properties: &armnetwork.TroubleshootingProperties{
			StorageID:   &stroageAountId,
			StoragePath: to.StringPtr("https://" + storageAccountName + ".blob.core.windows.net/troubleshooting"),
		},
		TargetResourceID: &gatewayId,
	}
	err = GetNetworkWatcherTroubleshooting(ctx, networkWatcherName, troubleshootingParameters)
	if err != nil {
		t.Fatalf("failed to get network watcher troubleshooting: % +v", err)
	}
	t.Logf("got network watcher troubleshooting")

	queryTroubleshootingParameters := armnetwork.QueryTroubleshootingParameters{
		TargetResourceID: &gatewayId,
	}
	err = GetNetworkWatcherTroubleshootingResult(ctx, networkWatcherName, queryTroubleshootingParameters)
	if err != nil {
		t.Fatalf("failed to get network watcher troubleshooting result: % +v", err)
	}
	t.Logf("got network watcher troubleshooting result")

	verificationIPFlowParameters := armnetwork.VerificationIPFlowParameters{
		Direction:        armnetwork.DirectionOutbound.ToPtr(),
		LocalIPAddress:   ipConfigPro.PrivateIPAddress,
		LocalPort:        to.StringPtr("80"),
		Protocol:         armnetwork.IPFlowProtocolTCP.ToPtr(),
		RemoteIPAddress:  to.StringPtr("121.10.1.1"),
		RemotePort:       to.StringPtr("80"),
		TargetResourceID: &vmId,
	}
	err = VerifyNetworkWatcherIPFlow(ctx, networkWatcherName, verificationIPFlowParameters)
	if err != nil {
		t.Fatalf("failed to verify network watcher ip flow: %+v", err)
	}
	t.Logf("verified network watcher ip flow")

	flowLogStatusParameters := armnetwork.FlowLogStatusParameters{
		TargetResourceID: &securityGroupId,
	}
	err = GetNetworkWatcherFlowLogStatus(ctx, networkWatcherName, flowLogStatusParameters)
	if err != nil {
		t.Fatalf("failed to get network watcher flow log status: %+v", err)
	}
	t.Logf("got network watcher flow log status")

	flowLogInformationParameters := armnetwork.FlowLogInformation{
		Properties: &armnetwork.FlowLogProperties{
			Enabled:   to.BoolPtr(true),
			StorageID: &stroageAountId,
		},
		TargetResourceID: &securityGroupId,
	}
	err = SetNetworkWatcherFlowLogConfiguration(ctx, networkWatcherName, flowLogInformationParameters)
	if err != nil {
		t.Fatalf("failed to set network watcher Configure flow log: %+v", err)
	}
	t.Logf("set network watcher Configure flow log")

	networkConfigurationDiagnosticParameters := armnetwork.NetworkConfigurationDiagnosticParameters{
		Profiles: &[]*armnetwork.NetworkConfigurationDiagnosticProfile{{
			Destination:     to.StringPtr("12.11.12.14"),
			DestinationPort: to.StringPtr("12100"),
			Direction:       armnetwork.DirectionInbound.ToPtr(),
			Protocol:        to.StringPtr("TCP"),
			Source:          to.StringPtr("10.1.0.4"),
		}},
		TargetResourceID: &vmId,
	}
	err = GetNetworkConfigurationDiagnostic(ctx, networkWatcherName, networkConfigurationDiagnosticParameters)
	if err != nil {
		t.Fatalf("failed to get network configuration diagnostic: %+v", err)
	}
	t.Logf("got network configuration diagnostic")

	securityGroupViewParameters := armnetwork.SecurityGroupViewParameters{
		TargetResourceID: &vmId,
	}
	err = GetNetworkVMSecurityRules(ctx, networkWatcherName, securityGroupViewParameters)
	if err != nil {
		t.Fatalf("failed to get network security group view: %+v", err)
	}
	t.Logf("got network security group view")

	connectivityParameters := armnetwork.ConnectivityParameters{
		Destination: &armnetwork.ConnectivityDestination{
			Address: to.StringPtr("192.168.100.4"),
			Port:    new(int32),
		},
		PreferredIPVersion: armnetwork.IPVersionIPv4.ToPtr(),
		Source: &armnetwork.ConnectivitySource{
			Port:       to.Int32Ptr(3389),
			ResourceID: &vmId,
		},
	}
	err = CheckNetworkConnectivity(ctx, networkWatcherName, connectivityParameters)
	if err != nil {
		t.Fatalf("failed to check connectivity: %+v", err)
	}
	t.Logf("checked connectivity")

	topologyParameters := armnetwork.TopologyParameters{
		TargetResourceGroupName: to.StringPtr(config.GroupName()),
	}
	err = GetNetworkTopology(ctx, networkWatcherName, topologyParameters)
	if err != nil {
		t.Fatalf("failed to get Topology: %+v", err)
	}
	t.Logf("got Topology")

	err = ListNetworkWatcher(ctx)
	if err != nil {
		t.Fatalf("failed to list network watcher: %+v", err)
	}
	t.Logf("listed network watcher")

	err = ListAllNetworkWatcher(ctx)
	if err != nil {
		t.Fatalf("failed to list all network watcher: %+v", err)
	}
	t.Logf("listed all network watcher")

	err = GetNetworkWatcher(ctx, networkWatcherName)
	if err != nil {
		t.Fatalf("failed to get network watcher: %+v", err)
	}
	t.Logf("got network watcher")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateNetworkWatcherTags(ctx, networkWatcherName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for network watcher: %+v", err)
	}
	t.Logf("updated network watcher tags")

	err = DeleteNetworkWatcher(ctx, networkWatcherName)
	if err != nil {
		t.Fatalf("failed to delete network watcher: %+v", err)
	}
	t.Logf("deleted network watcher")
}
