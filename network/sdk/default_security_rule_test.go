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

func TestDefaultSecurityRule(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	networkSecurityGroupName := config.AppendRandomSuffix("networksecuritygroup")
	defaultSecurityRuleName := "AllowVnetInBound"

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

	err = GetDefaultSecurityRule(ctx, networkSecurityGroupName, defaultSecurityRuleName)
	if err != nil {
		t.Fatalf("failed to get default security rule: %+v", err)
	}
	t.Logf("got default security rule")

	err = ListDefaultSecurityRule(ctx, networkSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to list default security rule: %+v", err)
	}
	t.Logf("listed default security rule")
}
