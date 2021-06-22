// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
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
func CreateNetworkInterface(ctx context.Context, networkInterfaceName string, networkInterfaceParameters armnetwork.NetworkInterface) (string, *armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat, error) {
	client := getNetworkInterfacesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkInterfaceName,
		networkInterfaceParameters,
		nil,
	)

	if err != nil {
		return "", nil, err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", nil, err
	}

	id := *resp.NetworkInterface.ID

	var ipConfigProperties *armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat
	if len((resp.NetworkInterface.Properties.IPConfigurations)) > 0 {
		ipConfigProperties = (resp.NetworkInterface.Properties.IPConfigurations)[0].Properties
	}
	return id, ipConfigProperties, nil
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
func UpdateNetworkInterfaceTags(ctx context.Context, networkInterfaceName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getNetworkInterfacesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkInterfaceName,
		tagsObjectParameters,
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
