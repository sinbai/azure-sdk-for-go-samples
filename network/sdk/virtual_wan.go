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

func getVirtualWansClient() armnetwork.VirtualWansClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualWansClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualWans
func CreateVirtualWan(ctx context.Context, virtualWanName string) (string, error) {
	client := getVirtualWansClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualWanName,
		armnetwork.VirtualWAN{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
				Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
			},
			Properties: &armnetwork.VirtualWanProperties{
				DisableVPNEncryption: to.BoolPtr(false),
				Type:                 to.StringPtr("Basic"),
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}
	return poller.RawResponse.Request.URL.Path, nil

	// resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	// if err != nil {
	// 	return "", err
	// }
	// return *resp.VirtualWAN.ID, nil
}
