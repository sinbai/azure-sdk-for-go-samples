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

func getHubRouteTablesClient() armnetwork.HubRouteTablesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewHubRouteTablesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a RouteTable resource if it doesn't exist else updates the existing RouteTable
func CreateHubRouteTable(ctx context.Context, virtualHubName string, routeTableName string, routeTableParameters armnetwork.HubRouteTable) error {
	client := getHubRouteTablesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualHubName,
		routeTableName,
		routeTableParameters,
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

// Retrieves the details of a RouteTable
func GetHubRouteTable(ctx context.Context, virtualHubName string, routeTableName string) error {
	client := getHubRouteTablesClient()
	_, err := client.Get(ctx, config.GroupName(), virtualHubName, routeTableName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves the details of all RouteTables.
func ListHubRouteTable(ctx context.Context, virtualHubName string) error {
	client := getHubRouteTablesClient()
	pager := client.List(config.GroupName(), virtualHubName, nil)

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

// Deletes a RouteTable.
func DeleteHubRouteTable(ctx context.Context, virtualHubName string, routeTableName string) error {
	client := getHubRouteTablesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualHubName, routeTableName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
