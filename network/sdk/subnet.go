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

func getSubnetsClient() armnetwork.SubnetsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create SubNets
func CreateSubnet(ctx context.Context, virtualNetworkName string, subnetName string, subnetParameters armnetwork.Subnet) (string, error) {
	client := getSubnetsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkName,
		subnetName,
		subnetParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.Subnet.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.Subnet.ID, nil
}
