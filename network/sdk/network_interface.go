// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

var (
	ipConfigurationName string
	subnetName          string
)

func getNetworkInterfacesClient() armnetwork.NetworkInterfacesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkInterfacesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create NetworkInterfaces
func CreateNetworkInterface(ctx context.Context, networkInterfaceName string, publicIpAddressName string, virtualNetworkName string) (*string, error) {
	ipConfigurationName = config.AppendRandomSuffix("ipconfiguration")

	urlPathAddress := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{publicIpAddressName}"
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{publicIpAddressName}", url.PathEscape(publicIpAddressName))
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	urlPathSubNet := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkName}/subnets/{subnetName}"
	urlPathSubNet = strings.ReplaceAll(urlPathSubNet, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathSubNet = strings.ReplaceAll(urlPathSubNet, "{virtualNetworkName}", url.PathEscape(virtualNetworkName))
	urlPathSubNet = strings.ReplaceAll(urlPathSubNet, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))
	urlPathSubNet = strings.ReplaceAll(urlPathSubNet, "{subnetName}", url.PathEscape(subnetName))

	client := getNetworkInterfacesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkInterfaceName,
		armnetwork.NetworkInterface{
			Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
			Properties: &armnetwork.NetworkInterfacePropertiesFormat{
				EnableAcceleratedNetworking: to.BoolPtr(true),
				IPConfigurations: &[]armnetwork.NetworkInterfaceIPConfiguration{
					{
						Name: &ipConfigurationName,
						Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
							PublicIPAddress: &armnetwork.PublicIPAddress{
								Resource: armnetwork.Resource{
									ID: &urlPathAddress,
								},
							},
							Subnet: &armnetwork.Subnet{SubResource: armnetwork.SubResource{ID: &urlPathSubNet}},
						},
					},
				},
			},
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}

	return resp.NetworkInterface.ID, nil
}

// Gets the specified network interface in a specified resource group.
func GetNetworkInterface(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfacesClient()
	_, err := client.Get(ctx, config.GroupName(), networkInterfaceName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the network interface in a subscription.
func ListNetworkInterface(ctx context.Context) error {
	client := getNetworkInterfacesClient()
	pager := client.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Gets all the network interface in a subscription.
func ListAllNetworkInterface(ctx context.Context) error {
	client := getNetworkInterfacesClient()
	pager := client.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Updates network interface tags.
func UpdateNetworkInterfaceTags(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfacesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkInterfaceName,
		armnetwork.TagsObject{
			Tags: &map[string]string{"tag1": "value1", "tag2": "value2"},
		},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified network interface.
func DeleteNetworkInterface(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfacesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkInterfaceName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all route tables applied to a network interface
func BeginGetEffectiveRouteTable(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfacesClient()
	poller, err := client.BeginGetEffectiveRouteTable(ctx, config.GroupName(), networkInterfaceName, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all network security groups applied to a network interface.
func BeginListEffectiveRouteTable(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfacesClient()
	poller, err := client.BeginListEffectiveNetworkSecurityGroups(ctx, config.GroupName(), networkInterfaceName, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func getNetworkInterfaceIPConfigurationsClient() armnetwork.NetworkInterfaceIPConfigurationsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkInterfaceIPConfigurationsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Get all ip configurations in a network interface.
func ListNetworkInterfaceIpConfiguration(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfaceIPConfigurationsClient()
	pager := client.List(config.GroupName(), networkInterfaceName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getNetworkInterfaceLoadBalancersClient() armnetwork.NetworkInterfaceLoadBalancersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkInterfaceLoadBalancersClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// List all load balancers in a network interface.
func ListNetworkInterfaceLoadBalancer(ctx context.Context, networkInterfaceName string) error {
	client := getNetworkInterfaceLoadBalancersClient()
	pager := client.List(config.GroupName(), networkInterfaceName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getVirtualNetworksClient() armnetwork.VirtualNetworksClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworksClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a virtual network in the specified resource group
func CreateVirtualNetwork(ctx context.Context, virtualNetworkName string) error {
	client := getVirtualNetworksClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},

			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/16"},
				},
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}

	err = CreateSubnet(ctx, virtualNetworkName, subnetName)
	if err != nil {
		return err
	}
	return nil
}

func getSubnetsClient() armnetwork.SubnetsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create SubNets
func CreateSubnet(ctx context.Context, virtualNetworkName string, subnetName string) error {
	client := getSubnetsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkName,
		subnetName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.0.0/24"),
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func getVirtualMachinesClient() armcompute.VirtualMachinesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Deletes the specified virtual machine.
func DeleteVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Create VirtualMachines
func CreateVirtualMachine(ctx context.Context, virtualMachineName string, nicId *string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualMachineName,
		armcompute.VirtualMachine{
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
								ID: nicId,
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
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
