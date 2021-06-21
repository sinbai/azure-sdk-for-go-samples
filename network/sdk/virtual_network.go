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

func getVirtualNetworksClient() armnetwork.VirtualNetworksClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworksClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a virtual network in the specified resource group
func CreateVirtualNetwork(ctx context.Context, virtualNetworkName string, virtualNetworkParameters armnetwork.VirtualNetwork) (string, error) {
	client := getVirtualNetworksClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkName,
		virtualNetworkParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.VirtualNetwork.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.VirtualNetwork.ID, nil
}

// Checks whether a private IP address is available for use.
func CheckIPAddressAvailability(ctx context.Context, virtualNetworkName string, ipAddress string) error {
	client := getVirtualNetworksClient()
	_, err := client.CheckIPAddressAvailability(ctx, config.GroupName(), virtualNetworkName, ipAddress, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists usage stats.
func ListUsageVirtualNetwork(ctx context.Context, virtualNetworkName string) error {
	client := getVirtualNetworksClient()
	pager := client.ListUsage(config.GroupName(), virtualNetworkName, nil)

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

// Gets the specified virtual network in a specified resource group.
func GetVirtualNetwork(ctx context.Context, virtualNetworkName string) error {
	client := getVirtualNetworksClient()
	_, err := client.Get(ctx, config.GroupName(), virtualNetworkName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the virtual network in a subscription.
func ListVirtualNetwork(ctx context.Context) error {
	client := getVirtualNetworksClient()
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

// Gets all the virtual network in a subscription.
func ListAllVirtualNetwork(ctx context.Context) error {
	client := getVirtualNetworksClient()
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

// Updates virtual network tags.
func UpdateVirtualNetworkTags(ctx context.Context, virtualNetworkName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getVirtualNetworksClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		virtualNetworkName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified virtual network.
func DeleteVirtualNetwork(ctx context.Context, virtualNetworkName string) error {
	client := getVirtualNetworksClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualNetworkName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
