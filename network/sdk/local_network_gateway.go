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

func getLocalNetworkGatewaysClient() armnetwork.LocalNetworkGatewaysClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLocalNetworkGatewaysClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create LocalNetworkGateways
func CreateLocalNetworkGateway(ctx context.Context, localNetworkGatewayName string, localNetworkGatewayParameters armnetwork.LocalNetworkGateway) (string, error) {
	client := getLocalNetworkGatewaysClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		localNetworkGatewayName,
		localNetworkGatewayParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.LocalNetworkGateway.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.LocalNetworkGateway.ID, nil
}

// Gets the specified local network gateway in a specified resource group.
func GetLocalNetworkGateway(ctx context.Context, localNetworkGatewayName string) error {
	client := getLocalNetworkGatewaysClient()
	_, err := client.Get(ctx, config.GroupName(), localNetworkGatewayName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the local network gateway in a subscription.
func ListLocalNetworkGateway(ctx context.Context) error {
	client := getLocalNetworkGatewaysClient()
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

// Updates local network gateway tags.
func UpdateLocalNetworkGatewayTags(ctx context.Context, localNetworkGatewayName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getLocalNetworkGatewaysClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		localNetworkGatewayName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified local network gateway.
func DeleteLocalNetworkGateway(ctx context.Context, localNetworkGatewayName string) error {
	client := getLocalNetworkGatewaysClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), localNetworkGatewayName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
