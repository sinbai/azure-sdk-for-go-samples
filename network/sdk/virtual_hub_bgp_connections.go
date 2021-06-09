// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getVirtualHubBgpConnectionsClient() armnetwork.VirtualHubBgpConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualHubBgpConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Retrieves the details of all VirtualHubBgpConnections.
func ListVirtualHubBgpConnection(ctx context.Context, virtualHubName string) error {
	client := getVirtualHubBgpConnectionsClient()
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

// Retrieves a list of routes the virtual hub bgp connection is advertising to the specified peer.
func ListVirtualHubBgpConnectionAdvertisedRoute(ctx context.Context, hubName string, connectionName string) error {
	client := getVirtualHubBgpConnectionsClient()
	_, err := client.BeginListAdvertisedRoutes(ctx, config.GroupName(), hubName, connectionName, nil)

	if err != nil {
		return err
	}
	return nil
}

// Retrieves a list of routes the virtual hub bgp connection has learned.
func ListVirtualHubBgpConnectionLearnedRoutes(ctx context.Context, hubName string, connectionName string) error {
	client := getVirtualHubBgpConnectionsClient()
	_, err := client.BeginListLearnedRoutes(ctx, config.GroupName(), hubName, connectionName, nil)

	if err != nil {
		return err
	}
	return nil
}
