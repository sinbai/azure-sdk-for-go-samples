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

func getVirtualWansClient() armnetwork.VirtualWansClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualWansClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualWans
func CreateVirtualWan(ctx context.Context, virtualWanName string, virtualWANParameters armnetwork.VirtualWAN) (string, error) {
	client := getVirtualWansClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualWanName,
		virtualWANParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.VirtualWAN.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.VirtualWAN.ID, nil
}

// Gets the specified virtual wan in a specified resource group.
func GetVirtualWan(ctx context.Context, virtualWanName string) error {
	client := getVirtualWansClient()
	_, err := client.Get(ctx, config.GroupName(), virtualWanName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the virtual wan in a subscription.
func ListVirtualWan(ctx context.Context) error {
	client := getVirtualWansClient()
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

// Updates virtual wan tags.
func UpdateVirtualWanTags(ctx context.Context, virtualWanName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getVirtualWansClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		virtualWanName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified virtual wan.
func DeleteVirtualWan(ctx context.Context, virtualWanName string) error {
	client := getVirtualWansClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualWanName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all virtual wan in a resource group.
func ListVirtualWanByResourceGroup(ctx context.Context) error {
	client := getVirtualWansClient()
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
