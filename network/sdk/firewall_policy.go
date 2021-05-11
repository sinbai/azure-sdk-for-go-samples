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

func getFirewallPolicysClient() armnetwork.FirewallPoliciesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewFirewallPoliciesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

//  Creates or updates the specified Firewall Policy.
func CreateFirewallPolicy(ctx context.Context, firewallPolicyName string) error {
	client := getFirewallPolicysClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		firewallPolicyName,
		armnetwork.FirewallPolicy{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
				Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
			},
			Properties: &armnetwork.FirewallPolicyPropertiesFormat{
				ThreatIntelMode: armnetwork.AzureFirewallThreatIntelModeAlert.ToPtr(),
			},
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

// Gets the specified firewall policy in a specified resource group.
func GetFirewallPolicy(ctx context.Context, firewallPolicyName string) error {
	client := getFirewallPolicysClient()
	_, err := client.Get(ctx, config.GroupName(), firewallPolicyName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the firewall policy in a subscription.
func ListFirewallPolicy(ctx context.Context) error {
	client := getFirewallPolicysClient()
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

// Gets all the firewall policy in a subscription.
func ListAllFirewallPolicy(ctx context.Context) error {
	client := getFirewallPolicysClient()
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

// Deletes the specified firewall policy.
func DeleteFirewallPolicy(ctx context.Context, firewallPolicyName string) error {
	client := getFirewallPolicysClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), firewallPolicyName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
