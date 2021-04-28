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

func getIPPrefixClient() armnetwork.PublicIPPrefixesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPPrefixesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create public IP prefix
func CreatePublicIPPrefix(ctx context.Context, prefixName string) error {
	client := getIPPrefixClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		prefixName,
		armnetwork.PublicIPPrefix{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(prefixName),
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.PublicIPPrefixPropertiesFormat{
				PrefixLength:           to.Int32Ptr(30),
				PublicIPAddressVersion: armnetwork.IPVersionIPv4.ToPtr(),
			},
			SKU: &armnetwork.PublicIPPrefixSKU{
				Name: armnetwork.PublicIPPrefixSKUNameStandard.ToPtr(),
			},
		},
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

// Gets the specified public IP prefix in a specified resource group.
func GetPublicIPPrefix(ctx context.Context, ipName string) error {
	client := getIPPrefixClient()
	_, err := client.Get(ctx, config.GroupName(), ipName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all public IP addresses in a resource group.
func ListPublicIPPrefix(ctx context.Context) error {
	client := getIPPrefixClient()
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
	client := getIPPrefixClient()
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
func UpdatePrefixTags(ctx context.Context, prefixName string) error {
	client := getIPPrefixClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		prefixName,
		armnetwork.TagsObject{
			Tags: &map[string]string{"tag1": "value1", "tag2": "value2"},
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
	client := getIPPrefixClient()
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
