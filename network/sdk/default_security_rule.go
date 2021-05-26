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

func getDefaultSecurityRulesClient() armnetwork.DefaultSecurityRulesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewDefaultSecurityRulesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Get the specified default network security rule.
func GetDefaultSecurityRule(ctx context.Context, networkSecurityGroupName string, defaultSecurityRuleName string) error {
	client := getDefaultSecurityRulesClient()
	_, err := client.Get(ctx, config.GroupName(), networkSecurityGroupName, defaultSecurityRuleName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all default security rules in a network security group.
func ListDefaultSecurityRule(ctx context.Context, networkSecurityGroupName string) error {
	client := getDefaultSecurityRulesClient()
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
