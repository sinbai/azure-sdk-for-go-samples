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

func getApplicationSecurityGroupsClient() armnetwork.ApplicationSecurityGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewApplicationSecurityGroupsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates an application security group.
func CreateApplicationSecurityGroup(ctx context.Context, applicationSecurityGroupName string, applicationSecurityGroupParameters armnetwork.ApplicationSecurityGroup) error {
	client := getApplicationSecurityGroupsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		applicationSecurityGroupName,
		applicationSecurityGroupParameters,
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

// Gets information about the specified application security group.
func GetApplicationSecurityGroup(ctx context.Context, applicationSecurityGroupName string) error {
	client := getApplicationSecurityGroupsClient()
	_, err := client.Get(ctx, config.GroupName(), applicationSecurityGroupName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the application security groups in a resource group.
func ListApplicationSecurityGroup(ctx context.Context) error {
	client := getApplicationSecurityGroupsClient()
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

// Gets all application security groups in a subscription.
func ListAllApplicationSecurityGroup(ctx context.Context) error {
	client := getApplicationSecurityGroupsClient()
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

// Updates an application security group's tags.
func UpdateApplicationSecurityGroupTags(ctx context.Context, applicationSecurityGroupName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getApplicationSecurityGroupsClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		applicationSecurityGroupName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified application security group.
func DeleteApplicationSecurityGroup(ctx context.Context, applicationSecurityGroupName string) error {
	client := getApplicationSecurityGroupsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), applicationSecurityGroupName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
