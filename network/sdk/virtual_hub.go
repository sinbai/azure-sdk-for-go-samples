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
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getVirtualHubsClient() armnetwork.VirtualHubsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualHubsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualHubs
func CreateVirtualHub(ctx context.Context, virtualHubName string, virtualWanID string) error {
	client := getVirtualHubsClient()

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualHubName,
		armnetwork.VirtualHub{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
				Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
			},
			Properties: &armnetwork.VirtualHubProperties{
				AddressPrefix: to.StringPtr("10.168.0.0/24"),
				SKU:           to.StringPtr("Basic"),
				VirtualWan: &armnetwork.SubResource{
					ID: &virtualWanID,
				},
			},
		},
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
