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

func getRouteTablesClient() armnetwork.RouteTablesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewRouteTablesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

func getRoutesClient() armnetwork.RoutesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewRoutesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create RouteTables
func CreateRouteTable(ctx context.Context, routeTableName string, routeName string) error {
	client := getRouteTablesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		routeTableName,
		armnetwork.RouteTable{
			Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
			// Properties: &armnetwork.RouteTablePropertiesFormat{
			// 	DisableBgpRoutePropagation: to.BoolPtr(true),
			// 	Routes: &[]armnetwork.Route{
			// 		{
			// 			Name: &routeName,
			// 			Properties: &armnetwork.RoutePropertiesFormat{
			// 				AddressPrefix: to.StringPtr("10.0.3.0/24"),
			// 				NextHopType:   armnetwork.RouteNextHopTypeVirtualNetworkGateway.ToPtr(),
			// 			},
			// 		},
			// 	},
			// },
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
func UpdateRouteTableTags(ctx context.Context, routeTableName string) error {
	client := getRouteTablesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		routeTableName,
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

// Creates or updates a route in the specified route table.
func CreateRoute(ctx context.Context, routeTableName string, routeName string) error {
	client := getRoutesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		routeTableName,
		routeName,
		armnetwork.Route{
			Properties: &armnetwork.RoutePropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.3.0/24"),
				NextHopType:   armnetwork.RouteNextHopTypeVirtualNetworkGateway.ToPtr(),
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

// Gets the specified route from a route table.
func GetRoute(ctx context.Context, routeTableName string, routeName string) error {
	client := getRoutesClient()
	_, err := client.Get(ctx, config.GroupName(), routeTableName, routeName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all routes in a route table.
func ListRoute(ctx context.Context, routeTableName string) error {
	client := getRoutesClient()
	pager := client.List(config.GroupName(), routeTableName, nil)

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

// Deletes the specified route from a route table.
func DeleteRoute(ctx context.Context, routeTableName string, routeName string) error {
	client := getRoutesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), routeTableName, routeName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
