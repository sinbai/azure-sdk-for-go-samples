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

func getInboundSecurityRulesClient() armnetwork.InboundSecurityRuleClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewInboundSecurityRuleClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates the specified Network Virtual Appliance Inbound Security Rules
func CreateInboundSecurityRule(ctx context.Context, networkVirtualApplianceName string, ruleCollectionName string, inboundSecurityRuleParameters armnetwork.InboundSecurityRule) error {
	client := getInboundSecurityRulesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkVirtualApplianceName,
		ruleCollectionName,
		inboundSecurityRuleParameters,
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
