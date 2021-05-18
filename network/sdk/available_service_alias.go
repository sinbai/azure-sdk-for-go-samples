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

func getAvailableServiceAliasesClient() armnetwork.AvailableServiceAliasesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewAvailableServiceAliasesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets all available service aliases for this subscription in this region.
func ListAvailableServiceAlias(ctx context.Context) error {
	client := getAvailableServiceAliasesClient()
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

// Gets all available service aliases for this resource group in this region.
func ListAvailableServiceAliasByResourceGroup(ctx context.Context) error {
	client := getAvailableServiceAliasesClient()
	pager := client.ListByResourceGroup(config.GroupName(), config.Location(), nil)
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
