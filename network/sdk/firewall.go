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
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getFirewallsClient() armnetwork.AzureFirewallsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewAzureFirewallsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create Firewalls
func CreateFirewall(ctx context.Context, firewallName string, azureFirewallParameters armnetwork.AzureFirewall) error {
	client := getFirewallsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		firewallName,
		azureFirewallParameters,
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 120*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified firewall in a specified resource group.
func GetFirewall(ctx context.Context, firewallName string) error {
	client := getFirewallsClient()
	_, err := client.Get(ctx, config.GroupName(), firewallName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the firewall in a subscription.
func ListFirewall(ctx context.Context) error {
	client := getFirewallsClient()
	pager := client.List(config.GroupName(), nil)

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

// Gets all the firewall in a subscription.
func ListAllFirewall(ctx context.Context) error {
	client := getFirewallsClient()
	pager := client.ListAll(nil)
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

// Updates firewall tags.
func UpdateFirewallTags(ctx context.Context, firewallName string) error {
	client := getFirewallsClient()
	poller, err := client.BeginUpdateTags(
		ctx,
		config.GroupName(),
		firewallName,
		armnetwork.TagsObject{
			Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
		},
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

// Deletes the specified firewall.
func DeleteFirewall(ctx context.Context, firewallName string) error {
	client := getFirewallsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), firewallName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func getAzureFirewallFqdnTagsClient() armnetwork.AzureFirewallFqdnTagsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewAzureFirewallFqdnTagsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets all the azure firewall fqdn tag in a subscription.
func ListAllAzureFirewallFqdnTag(ctx context.Context) error {
	client := getAzureFirewallFqdnTagsClient()
	pager := client.ListAll(nil)
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
