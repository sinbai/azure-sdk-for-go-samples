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

func getPublicIPAddressClient() armnetwork.PublicIPAddressesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPAddressesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create public IP address
func CreatePublicIPAddress(ctx context.Context, addressName string) error {
	client := getPublicIPAddressClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		addressName,
		armnetwork.PublicIPAddress{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(addressName),
				Location: to.StringPtr(config.Location()),
			},

			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
				PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
			},
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: armnetwork.PublicIPAddressSKUNameStandard.ToPtr(),
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

// Gets the specified public IP address in a specified resource group.
func GetPublicIPAddress(ctx context.Context, addressName string) error {
	client := getPublicIPAddressClient()
	_, err := client.Get(ctx, config.GroupName(), addressName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the public IP prefixes in a subscription.
func ListPublicIPAddress(ctx context.Context) error {
	client := getPublicIPAddressClient()
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

// Gets all the public IP addresses in a subscription.
func ListAllPublicIPAddress(ctx context.Context) error {
	client := getPublicIPAddressClient()
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

// Updates public IP address tags.
func UpdatePublicIPAddressTags(ctx context.Context, addressName string) error {
	client := getPublicIPAddressClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		addressName,
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

// Deletes the specified public IP address.
func DeletePublicIPAddress(ctx context.Context, addressName string) error {
	client := getPublicIPAddressClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), addressName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
