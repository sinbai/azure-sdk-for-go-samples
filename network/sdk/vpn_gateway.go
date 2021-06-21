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

func getVpnGatewaysClient() armnetwork.VPNGatewaysClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNGatewaysClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a virtual wan vpn gateway if it doesn't exist else updates the existing gateway
func CreateVpnGateway(ctx context.Context, vpnGatewayName string, vpnGatewayParameters armnetwork.VPNGateway) error {
	client := getVpnGatewaysClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		vpnGatewayName,
		vpnGatewayParameters,
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

// Retrieves the details of a virtual wan vpn gateway.
func GetVpnGateway(ctx context.Context, vpnGatewayName string) error {
	client := getVpnGatewaysClient()
	_, err := client.Get(ctx, config.GroupName(), vpnGatewayName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all the VpnGateways in a subscription.
func ListVpnGateway(ctx context.Context) error {
	client := getVpnGatewaysClient()
	pager := client.List(nil)

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

// Deletes a virtual wan vpn gatewa.
func DeleteVpnGateway(ctx context.Context, vpnGatewayName string) error {
	client := getVpnGatewaysClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), vpnGatewayName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all vpn gateway in a resource group.
func ListVpnGatewayByResourceGroup(ctx context.Context) error {
	client := getVpnGatewaysClient()
	pager := client.ListByResourceGroup(config.GroupName(), nil)
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

// Resets the primary of the vpn gateway in the specified resource group.
func ResetVpnGateway(ctx context.Context, vpnGatewayName string) error {
	client := getVpnGatewaysClient()
	poller, err := client.BeginReset(ctx, config.GroupName(), vpnGatewayName, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
