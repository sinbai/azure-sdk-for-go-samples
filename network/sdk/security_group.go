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

func getNetworkSecurityGroupsClient() armnetwork.NetworkSecurityGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create NetworkSecurityGroups
func CreateNetworkSecurityGroup(ctx context.Context, networkSecurityGroupName string, networkSecurityGroupParameters armnetwork.NetworkSecurityGroup) (string, error) {
	client := getNetworkSecurityGroupsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkSecurityGroupName,
		networkSecurityGroupParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.NetworkSecurityGroup.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.NetworkSecurityGroup.ID, nil
}

// Gets the specified network security group in a specified resource group.
func GetNetworkSecurityGroup(ctx context.Context, networkSecurityGroupName string) error {
	client := getNetworkSecurityGroupsClient()
	_, err := client.Get(ctx, config.GroupName(), networkSecurityGroupName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the network security group in a subscription.
func ListNetworkSecurityGroup(ctx context.Context) error {
	client := getNetworkSecurityGroupsClient()
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

// Gets all the network security group in a subscription.
func ListAllNetworkSecurityGroup(ctx context.Context) error {
	client := getNetworkSecurityGroupsClient()
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

// Updates network security group tags.
func UpdateNetworkSecurityGroupTags(ctx context.Context, networkSecurityGroupName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getNetworkSecurityGroupsClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkSecurityGroupName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified network security group.
func DeleteNetworkSecurityGroup(ctx context.Context, networkSecurityGroupName string) error {
	client := getNetworkSecurityGroupsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkSecurityGroupName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
