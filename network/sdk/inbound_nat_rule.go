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

func getInboundNatRulesClient() armnetwork.InboundNatRulesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewInboundNatRulesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a load balancer inbound nat rule.
func CreateInboundNatRule(ctx context.Context, loadBalancerName string, inboundNatRuleName string, inboundNatRuleParameters armnetwork.InboundNatRule) error {
	client := getInboundNatRulesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		loadBalancerName,
		inboundNatRuleName,
		inboundNatRuleParameters,
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

// Gets the specified load balancer inbound nat rule.
func GetInboundNatRule(ctx context.Context, loadBalancerName string, inboundNatRuleName string) error {
	client := getInboundNatRulesClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerName, inboundNatRuleName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the inbound nat rules in a load balancer.
func ListInboundNatRule(ctx context.Context, loadBalancerName string) error {
	client := getInboundNatRulesClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

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

// Deletes the specified load balancer inbound nat rule.
func DeleteInboundNatRule(ctx context.Context, loadBalancerName string, inboundNatRuleName string) error {
	client := getInboundNatRulesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), loadBalancerName, inboundNatRuleName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
