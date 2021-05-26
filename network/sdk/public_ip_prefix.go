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
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getPublicIPPrefixClient() armnetwork.PublicIPPrefixesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPPrefixesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create public IP prefix
func CreatePublicIPPrefix(ctx context.Context, prefixName string, publicIPPrefixParameters armnetwork.PublicIPPrefix) (string, error) {
	client := getPublicIPPrefixClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		prefixName,
		publicIPPrefixParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.PublicIPPrefix.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.PublicIPPrefix.ID, nil
}

// Gets the specified public IP prefix in a specified resource group.
func GetPublicIPPrefix(ctx context.Context, prefixName string) error {
	client := getPublicIPPrefixClient()
	_, err := client.Get(ctx, config.GroupName(), prefixName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all public IP prefix in a resource group.
func ListPublicIPPrefix(ctx context.Context) error {
	client := getPublicIPPrefixClient()
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

// Gets all the public IP prefix in a subscription.
func ListAllPublicIPPrefix(ctx context.Context) error {
	client := getPublicIPPrefixClient()
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

// Updates public IP prefix tags.
func UpdatePublicIPPrefixTags(ctx context.Context, prefixName string) error {
	client := getPublicIPPrefixClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		prefixName,
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

// Deletes the specified public IP prefix.
func DeletePublicIPPrefix(ctx context.Context, prefixName string) error {
	client := getPublicIPPrefixClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), prefixName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
