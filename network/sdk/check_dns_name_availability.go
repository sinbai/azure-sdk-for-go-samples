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

func getCheckDnsNameAvailabilitysClient() armnetwork.NetworkManagementClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkManagementClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets the specified check dns name availability in a specified resource group.
func GetCheckDnsNameAvailability(ctx context.Context, checkDnsNameAvailabilityName string) error {
	client := getCheckDnsNameAvailabilitysClient()
	_, err := client.CheckDNSNameAvailability(ctx, config.Location(), checkDnsNameAvailabilityName, nil)
	if err != nil {
		return err
	}
	return nil
}
