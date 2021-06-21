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

func getVirtualHubBgpConnectionClient() armnetwork.VirtualHubBgpConnectionClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualHubBgpConnectionClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a VirtualHubBgpConnection resource if it doesn't exist else updates the existing VirtualHubBgpConnection.
func CreateVirtualHubBgpConnection(ctx context.Context, virtualHubName string, connectionName string, bgpConnectionParameters armnetwork.BgpConnection) error {
	client := getVirtualHubBgpConnectionClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualHubName,
		connectionName,
		bgpConnectionParameters,
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

// Retrieves the details of a Virtual Hub Bgp Connection.
func GetVirtualHubBgpConnection(ctx context.Context, virtualHubName string, connectionName string) error {
	client := getVirtualHubBgpConnectionClient()
	_, err := client.Get(ctx, config.GroupName(), virtualHubName, connectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

//  Deletes a VirtualHubBgpConnection.
func DeleteVirtualHubBgpConnection(ctx context.Context, virtualHubName string, connectionName string) error {
	client := getVirtualHubBgpConnectionClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualHubName, connectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
