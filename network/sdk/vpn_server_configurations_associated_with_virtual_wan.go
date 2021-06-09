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

func getVpnServerConfigurationsAssociatedWithVirtualWansClient() armnetwork.VPNServerConfigurationsAssociatedWithVirtualWanClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNServerConfigurationsAssociatedWithVirtualWanClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gives the list of VpnServerConfigurations associated with Virtual Wan in a resource group.
func ListVpnServerConfigurationsAssociatedWithVirtualWan(ctx context.Context, virtualWANName string) error {
	client := getVpnServerConfigurationsAssociatedWithVirtualWansClient()
	poller, err := client.BeginList(ctx, config.GroupName(), virtualWANName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}

	return nil
}
