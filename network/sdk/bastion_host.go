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

func getBastionHostsClient() armnetwork.BastionHostsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewBastionHostsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create BastionHosts
func CreateBastionHost(ctx context.Context, bastionHostName string, bastionHostPro armnetwork.BastionHost) error {
	client := getBastionHostsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		bastionHostName,
		bastionHostPro,
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

// Gets the specified Bastion Host.
func GetBastionHost(ctx context.Context, bastionHostName string) error {
	client := getBastionHostsClient()
	_, err := client.Get(ctx, config.GroupName(), bastionHostName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all Bastion Hosts in a subscription.
func ListBastionHost(ctx context.Context) error {
	client := getBastionHostsClient()
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

// Deletes the specified Bastion Host.
func DeleteBastionHost(ctx context.Context, bastionHostName string) error {
	client := getBastionHostsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), bastionHostName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all bastion host in a resource group.
func ListBastionHostByResourceGroup(ctx context.Context) error {
	client := getBastionHostsClient()
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
