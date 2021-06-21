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

func getVirtualApplianceSkusClient() armnetwork.VirtualApplianceSKUsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualApplianceSKUsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Retrieves a single available sku for network virtual appliance.
func GetVirtualApplianceSku(ctx context.Context, virtualApplianceSkuName string) error {
	client := getVirtualApplianceSkusClient()
	_, err := client.Get(ctx, virtualApplianceSkuName, nil)
	if err != nil {
		return err
	}
	return nil
}

//  List all SKUs available for a virtual appliance.
func ListVirtualApplianceSku(ctx context.Context) error {
	client := getVirtualApplianceSkusClient()
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
