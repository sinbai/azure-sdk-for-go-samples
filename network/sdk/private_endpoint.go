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

func getPrivateEndpointsClient() armnetwork.PrivateEndpointsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPrivateEndpointsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create PrivateEndpoints
func CreatePrivateEndpoint(ctx context.Context, privateEndpointName string, privateEndpointPro armnetwork.PrivateEndpoint) error {
	client := getPrivateEndpointsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		privateEndpointName,
		privateEndpointPro,
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

// Gets the specified private endpoint in a specified resource group.
func GetPrivateEndpoint(ctx context.Context, privateEndpointName string) error {
	client := getPrivateEndpointsClient()
	_, err := client.Get(ctx, config.GroupName(), privateEndpointName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the private endpoint in a subscription.
func ListPrivateEndpoint(ctx context.Context) error {
	client := getPrivateEndpointsClient()
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

// Gets all the private endpoint in a subscription.
func ListAllPrivateEndpointBySubscription(ctx context.Context) error {
	client := getPrivateEndpointsClient()
	pager := client.ListBySubscription(nil)
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

// Deletes the specified private endpoint.
func DeletePrivateEndpoint(ctx context.Context, privateEndpointName string) error {
	client := getPrivateEndpointsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), privateEndpointName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
