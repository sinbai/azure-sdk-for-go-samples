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

func getRouteFilterClient() armnetwork.RouteFiltersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewRouteFiltersClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a route filter in a specified resource group.
func CreateRouteFilter(ctx context.Context, routeFilterName string, routeFilterParameters armnetwork.RouteFilter) error {
	client := getRouteFilterClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		routeFilterName,
		routeFilterParameters,
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

// Gets the specified route filter.
func GetRouteFilter(ctx context.Context, routeFilterName string) error {
	client := getRouteFilterClient()
	_, err := client.Get(ctx, config.GroupName(), routeFilterName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all route filters in a subscription.
func ListRouteFilter(ctx context.Context) error {
	client := getRouteFilterClient()
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

// Gets all route filters in a resource group.
func ListRouteFilterByResourceGroup(ctx context.Context) error {
	client := getRouteFilterClient()
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

// Updates route filter client tags.
func UpdateRouteFilterTags(ctx context.Context, routeFilterName string) error {
	client := getRouteFilterClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		routeFilterName,
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

// Deletes the specified route filter client.
func DeleteRouteFilter(ctx context.Context, routeFilterName string) error {
	client := getRouteFilterClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), routeFilterName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
