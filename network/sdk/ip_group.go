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

func getIPGroupClient() armnetwork.IPGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewIPGroupsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates an ipGroups in a specified resource group.
func CreateIPGroup(ctx context.Context, ipGroupName string, ipGroupPro armnetwork.IPGroup) error {
	client := getIPGroupClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		ipGroupName,
		ipGroupPro,
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

// Gets the specified ipGroups.
func GetIPGroup(ctx context.Context, ipGroupName string) error {
	client := getIPGroupClient()
	_, err := client.Get(ctx, config.GroupName(), ipGroupName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all IpGroups in a subscription.
func ListIPGroup(ctx context.Context) error {
	client := getIPGroupClient()
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

// Gets all IpGroups in a resource group.
func ListIPGroupByResourceGroup(ctx context.Context) error {
	client := getIPGroupClient()
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

// Deletes the specified ipGroups.
func DeleteIPGroup(ctx context.Context, ipGroupName string) error {
	client := getIPGroupClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), ipGroupName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
