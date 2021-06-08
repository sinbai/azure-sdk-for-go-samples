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
	"github.com/Azure/go-autorest/autorest/to"
)

func TestSecurityRule(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	securityRuleName := config.AppendRandomSuffix("securityrule")
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

	securityRuleParameters := armnetwork.SecurityRule{
		Properties: &armnetwork.SecurityRulePropertiesFormat{
			Access:                   armnetwork.SecurityRuleAccessDeny.ToPtr(),
			DestinationAddressPrefix: to.StringPtr("11.0.0.0/8"),
			DestinationPortRange:     to.StringPtr("8080"),
			Direction:                armnetwork.SecurityRuleDirectionOutbound.ToPtr(),
			Priority:                 to.Int32Ptr(100),
			Protocol:                 (*armnetwork.SecurityRuleProtocol)(to.StringPtr("*")),
			SourceAddressPrefix:      to.StringPtr("10.0.0.0/8"),
			SourcePortRange:          to.StringPtr("*"),
		},
	}
	err = CreateSecurityRule(ctx, networkSecurityGroupName, securityRuleName, securityRuleParameters)
	if err != nil {
		t.Fatalf("failed to create security rule: % +v", err)
	}
	t.Logf("created security rule")

	err = GetSecurityRule(ctx, networkSecurityGroupName, securityRuleName)
	if err != nil {
		t.Fatalf("failed to get security rule: %+v", err)
	}
	t.Logf("got security rule")

	err = ListSecurityRule(ctx, networkSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to list security rule: %+v", err)
	}
	t.Logf("listed security rule")

	err = DeleteSecurityRule(ctx, networkSecurityGroupName, securityRuleName)
	if err != nil {
		t.Fatalf("failed to delete security rule: %+v", err)
	}
	t.Logf("deleted security rule")

}
