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

func getFirewallPolicyRuleCollectionGroupsClient() armnetwork.FirewallPolicyRuleCollectionGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewFirewallPolicyRuleCollectionGroupsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates the specified FirewallPolicyRuleCollectionGroup.
func CreateFirewallPolicyRuleCollectionGroup(ctx context.Context, firewallPolicyName string, firewallPolicyRuleCollectionGroupName string, body string) error {
	client := getFirewallPolicyRuleCollectionGroupsClient()
	parameter := armnetwork.FirewallPolicyRuleCollectionGroupProperties{}
	parameter.UnmarshalJSON([]byte(body))

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		firewallPolicyName,
		firewallPolicyRuleCollectionGroupName,
		armnetwork.FirewallPolicyRuleCollectionGroup{
			Properties: &parameter,
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

// Gets the specified FirewallPolicyRuleCollectionGroup.
func GetFirewallPolicyRuleCollectionGroup(ctx context.Context, firewallPolicyName string, firewallPolicyRuleCollectionGroupName string) error {
	client := getFirewallPolicyRuleCollectionGroupsClient()
	_, err := client.Get(ctx, config.GroupName(), firewallPolicyName, firewallPolicyRuleCollectionGroupName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all FirewallPolicyRuleCollectionGroups in a FirewallPolicy resource.
func ListFirewallPolicyRuleCollectionGroup(ctx context.Context, firewallPolicyName string) error {
	client := getFirewallPolicyRuleCollectionGroupsClient()
	pager := client.List(config.GroupName(), firewallPolicyName, nil)

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

// Deletes the specified FirewallPolicyRuleCollectionGroup.
func DeleteFirewallPolicyRuleCollectionGroup(ctx context.Context, firewallPolicyName string, firewallPolicyRuleCollectionGroupName string) error {
	client := getFirewallPolicyRuleCollectionGroupsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), firewallPolicyName, firewallPolicyRuleCollectionGroupName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
