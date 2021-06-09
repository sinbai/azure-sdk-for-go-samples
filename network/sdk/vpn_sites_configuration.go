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

func getVpnSitesConfigurationsClient() armnetwork.VPNSitesConfigurationClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNSitesConfigurationClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gives the sas-url to download the configurations for vpn-sites in a resource group.
func DownloadVpnSitesConfiguration(ctx context.Context, virtualWANName string, request armnetwork.GetVPNSitesConfigurationRequest) error {
	client := getVpnSitesConfigurationsClient()
	poller, err := client.BeginDownload(ctx, config.GroupName(), virtualWANName, request, nil)

	if err != nil {
		return nil
	}

	_, err = poller.PollUntilDone(ctx, 300*time.Second)
	if err != nil {
		return err
	}

	return nil
}
