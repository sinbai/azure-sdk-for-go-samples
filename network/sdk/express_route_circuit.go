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

func getExpressRouteCircuitsClient() armnetwork.ExpressRouteCircuitsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewExpressRouteCircuitsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create ExpressRouteCircuits
func CreateExpressRouteCircuit(ctx context.Context, expressRouteCircuitName string, expressRouteCircuitParameters armnetwork.ExpressRouteCircuit) (string, error) {
	client := getExpressRouteCircuitsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		expressRouteCircuitName,
		expressRouteCircuitParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	return *resp.ExpressRouteCircuit.ID, nil
}

// Gets all stats from an express route circuit in a resource group.
func GetExpressRouteCircuitPeeringStats(ctx context.Context, expressRouteCircuitName string, peeringName string) error {
	client := getExpressRouteCircuitsClient()
	_, err := client.GetPeeringStats(ctx, config.GroupName(), expressRouteCircuitName, peeringName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the stats from an express route circuit in a resource group.
func GetExpressRouteCircuiStats(ctx context.Context, expressRouteCircuitName string) error {
	client := getExpressRouteCircuitsClient()
	_, err := client.GetStats(ctx, config.GroupName(), expressRouteCircuitName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified express route circuit in a specified resource group.
func GetExpressRouteCircuit(ctx context.Context, expressRouteCircuitName string) error {
	client := getExpressRouteCircuitsClient()
	_, err := client.Get(ctx, config.GroupName(), expressRouteCircuitName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the express route circuit in a subscription.
func ListExpressRouteCircuit(ctx context.Context) error {
	client := getExpressRouteCircuitsClient()
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

// Gets all the express route circuit in a subscription.
func ListAllExpressRouteCircuit(ctx context.Context) error {
	client := getExpressRouteCircuitsClient()
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

// Deletes the specified express route circuit.
func DeleteExpressRouteCircuit(ctx context.Context, expressRouteCircuitName string) error {
	client := getExpressRouteCircuitsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), expressRouteCircuitName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
