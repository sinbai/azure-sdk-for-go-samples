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

func getCustomIpPrefixesClient() armnetwork.CustomIPPrefixesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewCustomIPPrefixesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a custom IP prefix.
func CreateCustomIpPrefix(ctx context.Context, customIpPrefixName string, customIPPrefixParameters armnetwork.CustomIPPrefix) error {
	client := getCustomIpPrefixesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		customIpPrefixName,
		customIPPrefixParameters,
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

// Gets the specified custom IP prefix in a specified resource group.
func GetCustomIpPrefix(ctx context.Context, customIpPrefixName string) error {
	client := getCustomIpPrefixesClient()
	_, err := client.Get(ctx, config.GroupName(), customIpPrefixName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all custom IP prefixes in a resource group.
func ListCustomIpPrefix(ctx context.Context) error {
	client := getCustomIpPrefixesClient()
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

// Gets all the custom IP prefixes in a subscription.
func ListAllCustomIpPrefix(ctx context.Context) error {
	client := getCustomIpPrefixesClient()
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

// Updates custom ip prefix tags.
func UpdateCustomIpPrefixTags(ctx context.Context, customIpPrefixName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getCustomIpPrefixesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		customIpPrefixName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified custom ip prefix.
func DeleteCustomIpPrefix(ctx context.Context, customIpPrefixName string) error {
	client := getCustomIpPrefixesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), customIpPrefixName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
