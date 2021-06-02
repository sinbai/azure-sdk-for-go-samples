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

func getVpnSitesClient() armnetwork.VPNSitesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNSitesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a VpnSite resource if it doesn't exist else updates the existing VpnSite.
func CreateVpnSite(ctx context.Context, vpnSiteName string, vpnSiteParameters armnetwork.VPNSite) (string, error) {
	client := getVpnSitesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		vpnSiteName,
		vpnSiteParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}
	if resp.VPNSite.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.VPNSite.ID, nil
}

// Lists all the VpnSites in a subscription.
func ListVpnSite(ctx context.Context) error {
	client := getVpnSitesClient()
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

// Updates vpn site tags.
func UpdateVpnSiteTags(ctx context.Context, vpnSiteName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getVpnSitesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		vpnSiteName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes a VpnSite
func DeleteVpnSite(ctx context.Context, vpnSiteName string) error {
	client := getVpnSitesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), vpnSiteName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Lists all the vpnSites in a resource group.
func ListVpnSiteByResourceGroup(ctx context.Context) error {
	client := getVpnSitesClient()
	pager := client.ListByResourceGroup(config.GroupName(), nil)
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
