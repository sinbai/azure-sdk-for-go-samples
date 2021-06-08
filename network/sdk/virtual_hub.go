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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getVirtualHubsClient() armnetwork.VirtualHubsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualHubsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualHubs
func CreateVirtualHub(ctx context.Context, virtualHubName string, virtualHubParameters armnetwork.VirtualHub) (string, error) {
	client := getVirtualHubsClient()

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualHubName,
		virtualHubParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.VirtualHub.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.VirtualHub.ID, nil
}

// Retrieves the details of a VirtualHub
func GetVirtualHub(ctx context.Context, virtualHubName string) error {
	client := getVirtualHubsClient()
	_, err := client.Get(ctx, config.GroupName(), virtualHubName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all the VirtualHubs in a subscription.
func ListVirtualHub(ctx context.Context) error {
	client := getVirtualHubsClient()
	pager := client.List(nil)

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

//  Updates VirtualHub tags.
func UpdateVirtualHubTags(ctx context.Context, virtualHubName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getVirtualHubsClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		virtualHubName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes a VirtualHub.
func DeleteVirtualHub(ctx context.Context, virtualHubName string) error {
	client := getVirtualHubsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualHubName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Lists all the VirtualHubs in a resource group.
func ListVirtualHubByResourceGroup(ctx context.Context) error {
	client := getVirtualHubsClient()
	pager := client.ListByResourceGroup(config.GroupName(), nil)
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
