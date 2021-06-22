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

func getVirtualNetworkGatewaysClient() armnetwork.VirtualNetworkGatewaysClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworkGatewaysClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a virtual network gateway in the specified resource group.
func CreateVirtualNetworkGateway(ctx context.Context, virtualNetworkGatewayName string, virtualNetworkGateway armnetwork.VirtualNetworkGateway) (string, error) {
	client := getVirtualNetworkGatewaysClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkGatewayName,
		virtualNetworkGateway,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	return *resp.VirtualNetworkGateway.ID, nil
}

// Gets all the connections in a virtual network gateway.
func ListVirtualNetworkGatewayConnections(ctx context.Context, virtualNetworkGatewayName string) error {
	client := getVirtualNetworkGatewaysClient()
	pager := client.ListConnections(config.GroupName(), virtualNetworkGatewayName, nil)

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

// Gets the specified virtual network gateway in a specified resource group.
func GetVirtualNetworkGateway(ctx context.Context, virtualNetworkGatewayName string) error {
	client := getVirtualNetworkGatewaysClient()
	_, err := client.Get(ctx, config.GroupName(), virtualNetworkGatewayName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all virtual network gateways by resource group.
func ListVirtualNetworkGateway(ctx context.Context) error {
	client := getVirtualNetworkGatewaysClient()
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

// This operation retrieves a list of routes the virtual network gateway is advertising to the specified peer.
func BeginGetVirtualNetworkGatewayAdvertisedRoute(ctx context.Context, virtualNetworkGatewayName string, peer string) error {
	client := getVirtualNetworkGatewaysClient()
	poller, err := client.BeginGetAdvertisedRoutes(ctx, config.GroupName(), virtualNetworkGatewayName, peer, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// The GetBgpPeerStatus operation retrieves the status of all BGP peers.
func BeginGetVirtualNetworkGatewayBgpPeerStatus(ctx context.Context, virtualNetworkGatewayName string) error {
	client := getVirtualNetworkGatewaysClient()
	poller, err := client.BeginGetBgpPeerStatus(ctx, config.GroupName(), virtualNetworkGatewayName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// This operation retrieves a list of routes the virtual network gateway has learned, including routes learned from BGP peers.
func BeginGetVirtualNetworkGatewayLearnedRoutes(ctx context.Context, virtualNetworkGatewayName string) error {
	client := getVirtualNetworkGatewaysClient()
	poller, err := client.BeginGetLearnedRoutes(ctx, config.GroupName(), virtualNetworkGatewayName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Resets the primary of the virtual network gateway in the specified resource group.
func BeginVirtualNetworkGatewayReset(ctx context.Context, virtualNetworkGatewayName string) error {
	client := getVirtualNetworkGatewaysClient()
	poller, err := client.BeginReset(ctx, config.GroupName(), virtualNetworkGatewayName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Updates a virtual network gateway tags.
func UpdateVirtualNetworkGatewayTags(ctx context.Context, virtualNetworkGatewayName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getVirtualNetworkGatewaysClient()
	poller, err := client.BeginUpdateTags(
		ctx,
		config.GroupName(),
		virtualNetworkGatewayName,
		tagsObjectParameters,
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

// Deletes the specified virtual network gateway.
func DeleteVirtualNetworkGateway(ctx context.Context, virtualNetworkGatewayName string) error {
	client := getVirtualNetworkGatewaysClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualNetworkGatewayName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
