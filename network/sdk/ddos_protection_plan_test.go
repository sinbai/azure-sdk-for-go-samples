// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestDdosProtectionPlan(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	ddosProtectionPlanName := config.AppendRandomSuffix("ddosprotectionplan")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateDdosProtectionPlan(ctx, ddosProtectionPlanName)
	if err != nil {
		t.Fatalf("failed to create ddos protection plan: % +v", err)
	}
	t.Logf("created ddos protection plan")

	err = GetDdosProtectionPlan(ctx, ddosProtectionPlanName)
	if err != nil {
		t.Fatalf("failed to get ddos protection plan: %+v", err)
	}
	t.Logf("got ddos protection plan")

	err = ListDdosProtectionPlan(ctx)
	if err != nil {
		t.Fatalf("failed to list ddos protection plan: %+v", err)
	}
	t.Logf("listed ddos protection plan")

	err = ListDdosProtectionPlanByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to listddos protection plan by resource group: %+v", err)
	}
	t.Logf("listedddos protection plan by resource group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateDdosProtectionPlanTags(ctx, ddosProtectionPlanName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for ddos protection plan: %+v", err)
	}
	t.Logf("updated ddos protection plan tags")

	err = DeleteDdosProtectionPlan(ctx, ddosProtectionPlanName)
	if err != nil {
		t.Fatalf("failed to delete ddos protection plan: %+v", err)
	}
	t.Logf("deleted ddos protection plan")

}
