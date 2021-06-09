// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getApplicationGatewayPrivateLinkResourcesClient() armnetwork.ApplicationGatewayPrivateLinkResourcesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewApplicationGatewayPrivateLinkResourcesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Lists all private link resources on an application gateway
func ListApplicationGatewayPrivateLinkResource(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewayPrivateLinkResourcesClient()
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
