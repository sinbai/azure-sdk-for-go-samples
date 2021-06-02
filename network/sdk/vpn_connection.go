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

func getVPNConnectionsClient() armnetwork.VPNConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a vpn connection to a scalable vpn gateway if it doesn't exist else updates the existing connection
func CreateVpnConnection(ctx context.Context, gatewayName string, connectionName string, vpnConnectionParameters armnetwork.VPNConnection) error {
	client := getVPNConnectionsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		gatewayName,
		connectionName,
		vpnConnectionParameters,
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

// Retrieves the details of a vpn connection.
func GetVpnConnection(ctx context.Context, gatewayName string, connectionName string) error {
	client := getVPNConnectionsClient()
	_, err := client.Get(ctx, config.GroupName(), gatewayName, connectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

//  Deletes a vpn connection.
func DeleteVpnConnection(ctx context.Context, gatewayName string, connectionName string) error {
	client := getVPNConnectionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), gatewayName, connectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves all vpn connections for a particular virtual wan vpn gateway.
func ListVpnConnectionByVpnGateway(ctx context.Context, gatewayName string) error {
	client := getVPNConnectionsClient()
	pager := client.ListByVPNGateway(config.GroupName(), gatewayName, nil)
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
