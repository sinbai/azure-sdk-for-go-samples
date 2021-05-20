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

func getVirtualNetworkPeeringsClient() armnetwork.VirtualNetworkPeeringsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworkPeeringsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualNetworkPeerings
func CreateVirtualNetworkPeering(ctx context.Context, virtualNetworkName string, virtualNetworkPeeringName string, virtualNetworkPeeringPro armnetwork.VirtualNetworkPeering) error {
	client := getVirtualNetworkPeeringsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkName,
		virtualNetworkPeeringName,
		virtualNetworkPeeringPro,
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

//  Gets the specified virtual network peering.
func GetVirtualNetworkPeering(ctx context.Context, virtualNetworkName string, virtualNetworkPeeringName string) error {
	client := getVirtualNetworkPeeringsClient()
	_, err := client.Get(ctx, config.GroupName(), virtualNetworkName, virtualNetworkPeeringName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all virtual network peerings in a virtual network.
func ListVirtualNetworkPeering(ctx context.Context, virtualNetworkName string) error {
	client := getVirtualNetworkPeeringsClient()
	pager := client.List(config.GroupName(), virtualNetworkName, nil)

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

// Deletes the specified virtual network peering.
func DeleteVirtualNetworkPeering(ctx context.Context, virtualNetworkName string, virtualNetworkPeeringName string) error {
	client := getVirtualNetworkPeeringsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualNetworkName, virtualNetworkPeeringName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
