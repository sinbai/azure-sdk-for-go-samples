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

func TestNetworkSecurityGroup(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	networkSecurityGroupName := config.AppendRandomSuffix("networksecuritygroup")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	networkSecurityGroupParameters := armnetwork.NetworkSecurityGroup{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
	}
	_, err = CreateNetworkSecurityGroup(ctx, networkSecurityGroupName, networkSecurityGroupParameters)
	if err != nil {
		t.Fatalf("failed to create network security group: % +v", err)
	}
	t.Logf("created network security group")

	err = GetNetworkSecurityGroup(ctx, networkSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to get network security group: %+v", err)
	}
	t.Logf("got network security group")

	err = ListNetworkSecurityGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list network security group: %+v", err)
	}
	t.Logf("listed network security group")

	err = ListAllNetworkSecurityGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list all network security group: %+v", err)
	}
	t.Logf("listed all network security group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateNetworkSecurityGroupTags(ctx, networkSecurityGroupName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for network security group: %+v", err)
	}
	t.Logf("updated network security group tags")

	err = DeleteNetworkSecurityGroup(ctx, networkSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to delete network security group: %+v", err)
	}
	t.Logf("deleted network security group")
}
