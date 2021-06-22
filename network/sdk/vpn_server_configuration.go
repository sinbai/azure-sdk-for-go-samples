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

func getVpnServerConfigurationsClient() armnetwork.VPNServerConfigurationsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVPNServerConfigurationsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a VpnServerConfiguration resource if it doesn't exist else updates the existing VpnServerConfiguration.
func CreateVpnServerConfiguration(ctx context.Context, vpnServerConfigurationName string, vpnServerConfigurationParameters armnetwork.VPNServerConfiguration) (string, error) {
	client := getVpnServerConfigurationsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		vpnServerConfigurationName,
		vpnServerConfigurationParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	return *resp.VPNServerConfiguration.ID, nil
}

// Gets the specified vpn server configuration in a specified resource group.
func GetVpnServerConfiguration(ctx context.Context, vpnServerConfigurationName string) error {
	client := getVpnServerConfigurationsClient()
	_, err := client.Get(ctx, config.GroupName(), vpnServerConfigurationName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the vpn server configuration in a subscription.
func ListVpnServerConfiguration(ctx context.Context) error {
	client := getVpnServerConfigurationsClient()
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

// Updates vpn server configuration tags.
func UpdateVpnServerConfigurationTags(ctx context.Context, vpnServerConfigurationName string, vpnServerConfigurationParameters armnetwork.TagsObject) error {
	client := getVpnServerConfigurationsClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		vpnServerConfigurationName,
		vpnServerConfigurationParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified vpn server configuration.
func DeleteVpnServerConfiguration(ctx context.Context, vpnServerConfigurationName string) error {
	client := getVpnServerConfigurationsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), vpnServerConfigurationName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all vpn server configuration in a resource group.
func ListVpnServerConfigurationByResourceGroup(ctx context.Context) error {
	client := getVpnServerConfigurationsClient()
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
