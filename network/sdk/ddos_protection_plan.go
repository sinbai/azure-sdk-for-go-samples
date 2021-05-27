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

func getDdosProtectionPlansClient() armnetwork.DdosProtectionPlansClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewDdosProtectionPlansClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create DdosProtectionPlans
func CreateDdosProtectionPlan(ctx context.Context, ddosProtectionPlanName string) error {
	client := getDdosProtectionPlansClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		ddosProtectionPlanName,
		armnetwork.DdosProtectionPlan{
			Location: to.StringPtr(config.Location()),
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

// Gets the specified ddos protection plan in a specified resource group.
func GetDdosProtectionPlan(ctx context.Context, ddosProtectionPlanName string) error {
	client := getDdosProtectionPlansClient()
	_, err := client.Get(ctx, config.GroupName(), ddosProtectionPlanName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the ddos protection plan in a subscription.
func ListDdosProtectionPlan(ctx context.Context) error {
	client := getDdosProtectionPlansClient()
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

// Updates ddos protection plan tags.
func UpdateDdosProtectionPlanTags(ctx context.Context, ddosProtectionPlanName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getDdosProtectionPlansClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		ddosProtectionPlanName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified ddos protection plan.
func DeleteDdosProtectionPlan(ctx context.Context, ddosProtectionPlanName string) error {
	client := getDdosProtectionPlansClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), ddosProtectionPlanName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all ddos protection plan in a resource group.
func ListDdosProtectionPlanByResourceGroup(ctx context.Context) error {
	client := getDdosProtectionPlansClient()
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
