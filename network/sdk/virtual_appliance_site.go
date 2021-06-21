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
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
)

func getVirtualApplianceSitesClient() armnetwork.VirtualApplianceSitesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualApplianceSitesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates the specified Network Virtual Appliance Site.
func CreateVirtualApplianceSite(ctx context.Context, networkVirtualApplianceName string, siteName string, virtualApplianceSiteParameters armnetwork.VirtualApplianceSite) error {
	client := getVirtualApplianceSitesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkVirtualApplianceName,
		siteName,
		virtualApplianceSiteParameters,
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified Virtual Appliance Site.
func GetVirtualApplianceSite(ctx context.Context, networkVirtualApplianceName string, siteName string) error {
	client := getVirtualApplianceSitesClient()
	_, err := client.Get(ctx, config.GroupName(), networkVirtualApplianceName, siteName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all Network Virtual Appliance Sites in a Network Virtual Appliance resource.
func ListVirtualApplianceSite(ctx context.Context, networkVirtualApplianceName string) error {
	client := getVirtualApplianceSitesClient()
	pager := client.List(config.GroupName(), networkVirtualApplianceName, nil)

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

// Deletes the specified site from a Virtual Appliance.
func DeleteVirtualApplianceSite(ctx context.Context, networkVirtualApplianceName string, siteName string) error {
	client := getVirtualApplianceSitesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkVirtualApplianceName, siteName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
