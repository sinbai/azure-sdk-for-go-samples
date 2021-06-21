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

func getNetworkVirtualAppliancesClient() armnetwork.NetworkVirtualAppliancesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkVirtualAppliancesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates the specified Network Virtual Appliance.
func CreateNetworkVirtualAppliance(ctx context.Context, networkVirtualApplianceName string, parameters armnetwork.NetworkVirtualAppliance) error {
	client := getNetworkVirtualAppliancesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkVirtualApplianceName,
		parameters,
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

// Gets the specified Network Virtual Appliance.
func GetNetworkVirtualAppliance(ctx context.Context, networkVirtualApplianceName string) error {
	client := getNetworkVirtualAppliancesClient()
	_, err := client.Get(ctx, config.GroupName(), networkVirtualApplianceName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all Network Virtual Appliances in a subscription.
func ListNetworkVirtualAppliance(ctx context.Context) error {
	client := getNetworkVirtualAppliancesClient()
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

// Updates a Network Virtual Appliance.
func UpdateNetworkVirtualApplianceTags(ctx context.Context, networkVirtualApplianceName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getNetworkVirtualAppliancesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkVirtualApplianceName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

//  Deletes the specified Network Virtual Appliance.
func DeleteNetworkVirtualAppliance(ctx context.Context, networkVirtualApplianceName string) error {
	client := getNetworkVirtualAppliancesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkVirtualApplianceName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Lists all Network Virtual Appliances in a resource group.
func ListNetworkVirtualApplianceByResourceGroup(ctx context.Context) error {
	client := getNetworkVirtualAppliancesClient()
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
