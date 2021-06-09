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

func getApplicationGatewayPrivateEndpointConnectionsClient() armnetwork.ApplicationGatewayPrivateEndpointConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewApplicationGatewayPrivateEndpointConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Updates the specified private endpoint connection on application gateway
func UpdateApplicationGatewayPrivateEndpointConnection(ctx context.Context, applicationGatewayName string, connectionName string,
	parameters armnetwork.ApplicationGatewayPrivateEndpointConnection) error {
	client := getApplicationGatewayPrivateEndpointConnectionsClient()
	poller, err := client.BeginUpdate(
		ctx,
		config.GroupName(),
		applicationGatewayName,
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

// Gets the specified private endpoint connection on application gateway.
func GetApplicationGatewayPrivateEndpointConnection(ctx context.Context, applicationGatewayName string, connectionName string) error {
	client := getApplicationGatewayPrivateEndpointConnectionsClient()
	_, err := client.Get(ctx, config.GroupName(), applicationGatewayName, connectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

//  Lists all private endpoint connections on an application gateway.
func ListApplicationGatewayPrivateEndpointConnection(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewayPrivateEndpointConnectionsClient()
	pager := client.List(config.GroupName(), applicationGatewayName, nil)

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

// Deletes the specified private endpoint connection on application gateway.
func DeleteApplicationGatewayPrivateEndpointConnection(ctx context.Context, applicationGatewayName string, connectionName string) error {
	client := getApplicationGatewayPrivateEndpointConnectionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), applicationGatewayName, connectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
