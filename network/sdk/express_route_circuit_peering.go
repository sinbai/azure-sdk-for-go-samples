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
	"github.com/Azure/go-autorest/autorest/to"
)

func getExpressRouteCircuitPeeringsClient() armnetwork.ExpressRouteCircuitPeeringsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewExpressRouteCircuitPeeringsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create ExpressRouteCircuitPeerings
func CreateExpressRouteCircuitPeering(ctx context.Context, circuitName string, expressRouteCircuitPeeringName string) error {
	client := getExpressRouteCircuitPeeringsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		circuitName,
		expressRouteCircuitPeeringName,
		armnetwork.ExpressRouteCircuitPeering{
			Properties: &armnetwork.ExpressRouteCircuitPeeringPropertiesFormat{
				PeerASN:                    to.Int64Ptr(10001),
				PrimaryPeerAddressPrefix:   to.StringPtr("102.0.0.0/30"),
				SecondaryPeerAddressPrefix: to.StringPtr("103.0.0.0/30"),
				VlanID:                     to.Int32Ptr(101),
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

// Gets the specified peering for the express route circuit.
func GetExpressRouteCircuitPeering(ctx context.Context, circuitName string, expressRouteCircuitPeeringName string) error {
	client := getExpressRouteCircuitPeeringsClient()
	_, err := client.Get(ctx, config.GroupName(), circuitName, expressRouteCircuitPeeringName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all peerings in a specified express route circuit.
func ListExpressRouteCircuitPeering(ctx context.Context, circuitName string) error {
	client := getExpressRouteCircuitPeeringsClient()
	pager := client.List(config.GroupName(), circuitName, nil)

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

// Deletes the specified peering from the specified express route circuit.
func DeleteExpressRouteCircuitPeering(ctx context.Context, circuitName string, expressRouteCircuitPeeringName string) error {
	client := getExpressRouteCircuitPeeringsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), circuitName, expressRouteCircuitPeeringName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}