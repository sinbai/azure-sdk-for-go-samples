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

func getWpFirewallPoliciesClient() armnetwork.WebApplicationFirewallPoliciesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewWebApplicationFirewallPoliciesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or update policy with specified rule set name within a resource group.
func CreateWebApplicationFirewallPolicy(ctx context.Context, firewallPolicyName string,
	webApplicationFirewallPolicyParameters armnetwork.WebApplicationFirewallPolicy) error {
	client := getWpFirewallPoliciesClient()
	_, err := client.CreateOrUpdate(
		ctx,
		config.GroupName(),
		firewallPolicyName,
		webApplicationFirewallPolicyParameters,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

// Gets the specified Firewall Policy.
func GetWebApplicationFirewallPolicy(ctx context.Context, firewallPolicyName string) error {
	client := getWpFirewallPoliciesClient()
	_, err := client.Get(ctx, config.GroupName(), firewallPolicyName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all Firewall Policies in a resource group.
func ListWebApplicationFirewallPolicy(ctx context.Context) error {
	client := getWpFirewallPoliciesClient()
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

// Gets all the Firewall Policies in a subscription.
func ListAllWebApplicationFirewallPolicy(ctx context.Context) error {
	client := getWpFirewallPoliciesClient()
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

// Deletes the specified Firewall Policy.
func DeleteWebApplicationFirewallPolicy(ctx context.Context, firewallPolicyName string) error {
	client := getWpFirewallPoliciesClient()
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
