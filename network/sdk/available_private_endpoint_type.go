// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
)

func getAvailablePrivateEndpointTypesClient() armnetwork.AvailablePrivateEndpointTypesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewAvailablePrivateEndpointTypesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets all the available private endpoint type in a subscription.
func ListAvailablePrivateEndpointType(ctx context.Context) error {
	client := getAvailablePrivateEndpointTypesClient()
	pager := client.List(config.Location(), nil)

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

// Gets all available private endpoint type in a resource group.
func ListAvailablePrivateEndpointTypeByResourceGroup(ctx context.Context) error {
	client := getAvailablePrivateEndpointTypesClient()
	pager := client.ListByResourceGroup(config.Location(), config.GroupName(), nil)
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
