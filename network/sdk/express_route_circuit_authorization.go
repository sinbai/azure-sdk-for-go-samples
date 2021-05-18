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

func getExpressRouteCircuitAuthorizationsClient() armnetwork.ExpressRouteCircuitAuthorizationsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewExpressRouteCircuitAuthorizationsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates an authorization in the specified express route circuit.
func CreateExpressRouteCircuitAuthorization(ctx context.Context, circuitName string, expressRouteCircuitAuthorizationName string) error {
	client := getExpressRouteCircuitAuthorizationsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		circuitName,
		expressRouteCircuitAuthorizationName,
		armnetwork.ExpressRouteCircuitAuthorization{},
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

// Gets the specified authorization from the specified express route circuit.
func GetExpressRouteCircuitAuthorization(ctx context.Context, circuitName string, expressRouteCircuitAuthorizationName string) error {
	client := getExpressRouteCircuitAuthorizationsClient()
	_, err := client.Get(ctx, config.GroupName(), circuitName, expressRouteCircuitAuthorizationName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all authorizations in an express route circuit.
func ListExpressRouteCircuitAuthorization(ctx context.Context, circuitName string) error {
	client := getExpressRouteCircuitAuthorizationsClient()
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

// Deletes the specified authorization from the specified express route circuit.
func DeleteExpressRouteCircuitAuthorization(ctx context.Context, circuitName string, expressRouteCircuitAuthorizationName string) error {
	client := getExpressRouteCircuitAuthorizationsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), circuitName, expressRouteCircuitAuthorizationName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
