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

func getExpressRouteCircuitConnectionsClient() armnetwork.ExpressRouteCircuitConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewExpressRouteCircuitConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a Express Route Circuit Connection in the specified express route circuits.
func CreateExpressRouteCircuitConnection(ctx context.Context, circuitName string, peeringName string, connectionName string, parameters armnetwork.ExpressRouteCircuitConnection) error {
	client := getExpressRouteCircuitConnectionsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		circuitName,
		peeringName,
		connectionName,
		parameters,
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

// Gets the specified Express Route Circuit Connection from the specified express route circuit.
func GetExpressRouteCircuitConnection(ctx context.Context, circuitName string, peeringName string, connectionName string) error {
	client := getExpressRouteCircuitConnectionsClient()
	_, err := client.Get(ctx, config.GroupName(), circuitName, peeringName, connectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all global reach connections associated with a private peering in an express route circuit.
func ListExpressRouteCircuitConnection(ctx context.Context, circuitName string, peeringName string) error {
	client := getExpressRouteCircuitConnectionsClient()
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

// Deletes the specified Express Route Circuit Connection from the specified express route circuit.
func DeleteExpressRouteCircuitConnection(ctx context.Context, circuitName string, peeringName string, connectionName string) error {
	client := getExpressRouteCircuitConnectionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), circuitName, peeringName, connectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
