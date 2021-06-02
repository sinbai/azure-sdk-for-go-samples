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

func getVpnSiteLinksClient() armnetwork.VPNSiteLinksClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNSiteLinksClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Lists all the vpnSiteLinks in a resource group for a vpn site.
func ListVpnSiteLinkByVPNSite(ctx context.Context, vpnSiteName string) error {
	client := getVpnSiteLinksClient()
	pager := client.ListByVPNSite(config.GroupName(), vpnSiteName, nil)

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

// Retrieves the details of a VPN site link
func GetVpnSiteLink(ctx context.Context, vpnSiteName string, vpnSiteLinkName string) error {
	client := getVpnSiteLinksClient()
	_, err := client.Get(ctx, config.GroupName(), vpnSiteName, vpnSiteLinkName, nil)

	if err != nil {
		return err
	}
	return nil
}
