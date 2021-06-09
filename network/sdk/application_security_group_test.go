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

func TestApplicationSecurityGroup(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	applicationSecurityGroupName := config.AppendRandomSuffix("applicationsecuritygroup")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	applicationSecurityGroupParameters := armnetwork.ApplicationSecurityGroup{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	err = CreateApplicationSecurityGroup(ctx, applicationSecurityGroupName, applicationSecurityGroupParameters)
	if err != nil {
		t.Fatalf("failed to create application security group: % +v", err)
	}
	t.Logf("created application security group")

	err = GetApplicationSecurityGroup(ctx, applicationSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to get application security group: %+v", err)
	}
	t.Logf("got application security group")

	err = ListApplicationSecurityGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list application security group: %+v", err)
	}
	t.Logf("listed application security group")

	err = ListAllApplicationSecurityGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list all application security group: %+v", err)
	}
	t.Logf("listed all application security group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateApplicationSecurityGroupTags(ctx, applicationSecurityGroupName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for application security group: %+v", err)
	}
	t.Logf("updated application security group tags")

	err = DeleteApplicationSecurityGroup(ctx, applicationSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to delete application security group: %+v", err)
	}
	t.Logf("deleted application security group")

}
