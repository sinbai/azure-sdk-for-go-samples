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

func TestServiceEndpointPolicy(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	config.SetLocation("eastus")

	serviceEndpointPolicyName := config.AppendRandomSuffix("serviceendpointpolicy")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.DefaultLocation())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateServiceEndpointPolicy(ctx, serviceEndpointPolicyName)
	if err != nil {
		t.Fatalf("failed to create service endpoint policy: % +v", err)
	}
	t.Logf("created service endpoint policy")

	err = GetServiceEndpointPolicy(ctx, serviceEndpointPolicyName)
	if err != nil {
		t.Fatalf("failed to get service endpoint policy: %+v", err)
	}
	t.Logf("got service endpoint policy")

	err = ListServiceEndpointPolicyByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to listservice endpoint policy by resource group: %+v", err)
	}
	t.Logf("listedservice endpoint policy by resource group")

	err = ListServiceEndpointPolicy(ctx)
	if err != nil {
		t.Fatalf("failed to list service endpoint policy: %+v", err)
	}
	t.Logf("listed service endpoint policy")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateServiceEndpointPolicyTags(ctx, serviceEndpointPolicyName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for service endpoint policy: %+v", err)
	}
	t.Logf("updated service endpoint policy tags")

	err = DeleteServiceEndpointPolicy(ctx, serviceEndpointPolicyName)
	if err != nil {
		t.Fatalf("failed to delete service endpoint policy: %+v", err)
	}
	t.Logf("deleted service endpoint policy")

}
