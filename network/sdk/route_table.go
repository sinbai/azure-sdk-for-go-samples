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
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getRouteTablesClient() armnetwork.RouteTablesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewRouteTablesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create RouteTables
func CreateRouteTable(ctx context.Context, routeTableName string) error {
	client := getRouteTablesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		routeTableName,
		armnetwork.RouteTable{
			Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
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

// Gets the specified route table in a specified resource group.
func GetRouteTable(ctx context.Context, routeTableName string) error {
	client := getRouteTablesClient()
	_, err := client.Get(ctx, config.GroupName(), routeTableName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the route table in a subscription.
func ListRouteTable(ctx context.Context) error {
	client := getRouteTablesClient()
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

// Gets all the route table in a subscription.
func ListAllRouteTable(ctx context.Context) error {
	client := getRouteTablesClient()
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

// Updates route table tags.
func UpdateRouteTableTags(ctx context.Context, routeTableName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getRouteTablesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		routeTableName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified route table.
func DeleteRouteTable(ctx context.Context, routeTableName string) error {
	client := getRouteTablesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), routeTableName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
