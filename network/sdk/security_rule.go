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

func getSecurityRulesClient() armnetwork.SecurityRulesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSecurityRulesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create SecurityRules
func CreateSecurityRule(ctx context.Context, networkSecurityGroupName string, securityRuleName string, securityRuleParameters armnetwork.SecurityRule) error {
	client := getSecurityRulesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkSecurityGroupName,
		securityRuleName,
		securityRuleParameters,
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

//  Get the specified network security rule.
func GetSecurityRule(ctx context.Context, networkSecurityGroupName string, securityRuleName string) error {
	client := getSecurityRulesClient()
	_, err := client.Get(ctx, config.GroupName(), networkSecurityGroupName, securityRuleName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all security rules in a network security group.
func ListSecurityRule(ctx context.Context, networkSecurityGroupName string) error {
	client := getSecurityRulesClient()
	pager := client.List(config.GroupName(), networkSecurityGroupName, nil)

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

// Deletes the specified network security rule.
func DeleteSecurityRule(ctx context.Context, networkSecurityGroupName string, securityRuleName string) error {
	client := getSecurityRulesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkSecurityGroupName, securityRuleName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
