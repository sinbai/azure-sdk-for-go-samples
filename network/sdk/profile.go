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

func getNetworkProfilesClient() armnetwork.NetworkProfilesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkProfilesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create NetworkProfiles
func CreateNetworkProfile(ctx context.Context, networkProfileName string, networkProfileParameters armnetwork.NetworkProfile) error {
	client := getNetworkProfilesClient()
	_, err := client.CreateOrUpdate(
		ctx,
		config.GroupName(),
		networkProfileName,
		networkProfileParameters,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

// Gets the specified network profile in a specified resource group.
func GetNetworkProfile(ctx context.Context, networkProfileName string) error {
	client := getNetworkProfilesClient()
	_, err := client.Get(ctx, config.GroupName(), networkProfileName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the network profile in a subscription.
func ListNetworkProfile(ctx context.Context) error {
	client := getNetworkProfilesClient()
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

// Gets all the network profile in a subscription.
func ListAllNetworkProfile(ctx context.Context) error {
	client := getNetworkProfilesClient()
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

// Updates network profile tags.
func UpdateNetworkProfileTags(ctx context.Context, networkProfileName string) error {
	client := getNetworkProfilesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkProfileName,
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

// Deletes the specified network profile.
func DeleteNetworkProfile(ctx context.Context, networkProfileName string) error {
	client := getNetworkProfilesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkProfileName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
