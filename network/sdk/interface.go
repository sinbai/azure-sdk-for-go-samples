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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

var (
	ipConfigurationName string
	//subnetName          string
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
func CreateNetworkInterface(ctx context.Context, networkInterfaceName string, publicIpAddressName string, virtualNetworkName string, subnetName string) (*string, error) {
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
				IPConfigurations: &[]*armnetwork.NetworkInterfaceIPConfiguration{
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

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}

	return &poller.RawResponse.Request.URL.Path, nil
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
			Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
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
