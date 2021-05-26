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

func getVirtualNetworkGatewayConnectionsClient() armnetwork.VirtualNetworkGatewayConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworkGatewayConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualNetworkGatewayConnections
func CreateVirtualNetworkGatewayConnection(ctx context.Context, virtualNetworkGatewayConnectionName string, gatewayConnectionParameters armnetwork.VirtualNetworkGatewayConnection) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkGatewayConnectionName,
		gatewayConnectionParameters,
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

// BeginSetSharedKey - The Put VirtualNetworkGatewayConnectionSharedKey operation sets the virtual network gateway connection shared key for passed virtual
// network gateway connection in the specified resource group through
// Network resource provider.
func BeginSetVirtualNetworkGatewayConnectionSharedKey(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	poller, err := client.BeginSetSharedKey(ctx, config.GroupName(), virtualNetworkGatewayConnectionName,
		armnetwork.ConnectionSharedKey{Value: to.StringPtr("AzureAbc124")}, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// GetSharedKey - The Get VirtualNetworkGatewayConnectionSharedKey operation retrieves information about the specified virtual network gateway connection
// shared key through Network resource provider.
func GetVirtualNetworkGatewayConnectionSharedKey(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.GetSharedKey(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified virtual network gateway connection in a specified resource group.
func GetVirtualNetworkGatewayConnection(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.Get(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the virtual network gateway connection in a subscription.
func ListVirtualNetworkGatewayConnection(ctx context.Context) error {
	client := getVirtualNetworkGatewayConnectionsClient()
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

// BeginResetSharedKey - The VirtualNetworkGatewayConnectionResetSharedKey operation resets the virtual network gateway connection shared key for passed
// virtual network gateway connection in the specified resource group
// through Network resource provider.
func BeginResetVirtualNetworkGatewayConnectionSharedKey(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	poller, err := client.BeginResetSharedKey(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, armnetwork.ConnectionResetSharedKey{
		KeyLength: to.Int32Ptr(128),
	}, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Updates virtual network gateway connection tags.
func UpdateVirtualNetworkGatewayConnectionTags(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	poller, err := client.BeginUpdateTags(
		ctx,
		config.GroupName(),
		virtualNetworkGatewayConnectionName,
		armnetwork.TagsObject{
			Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
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

// Deletes the specified virtual network gateway connection.
func DeleteVirtualNetworkGatewayConnection(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
