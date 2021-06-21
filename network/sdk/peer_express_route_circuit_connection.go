// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
)

func getPeerExpressRouteCircuitConnectionsClient() armnetwork.PeerExpressRouteCircuitConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPeerExpressRouteCircuitConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets the specified Peer Express Route Circuit Connection from the specified express route circuit.
func GetPeerExpressRouteCircuitConnection(ctx context.Context, circuitName string, peeringName string,
	connectionName string) error {
	client := getPeerExpressRouteCircuitConnectionsClient()
	_, err := client.Get(ctx, config.GroupName(), circuitName, peeringName,
		connectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all global reach peer connections associated with a private peering in an express route circuit.
func ListPeerExpressRouteCircuitConnection(ctx context.Context, circuitName string, peeringName string) error {
	client := getPeerExpressRouteCircuitConnectionsClient()
	pager := client.List(config.GroupName(), circuitName, peeringName, nil)

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
