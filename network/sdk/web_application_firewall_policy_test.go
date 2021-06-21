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
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestWebApplicationFirewallPolicy(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	firewallPolicyName := config.AppendRandomSuffix("wafpolicy")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	webApplicationFirewallPolicyParameters := armnetwork.WebApplicationFirewallPolicy{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.WebApplicationFirewallPolicyPropertiesFormat{
			ManagedRules: &armnetwork.ManagedRulesDefinition{
				ManagedRuleSets: []*armnetwork.ManagedRuleSet{
					{
						RuleSetType:    to.StringPtr("OWASP"),
						RuleSetVersion: to.StringPtr("3.0"),
					},
				},
			},
		},
	}
	err = CreateWebApplicationFirewallPolicy(ctx, firewallPolicyName, webApplicationFirewallPolicyParameters)
	if err != nil {
		t.Fatalf("failed to create web application firewall policy: %+v", err)
	}
	t.Logf("created web application firewall policy")

	err = GetWebApplicationFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to get web application firewall policy: %+v", err)
	}
	t.Logf("got web application firewall policy")

	err = ListWebApplicationFirewallPolicy(ctx)
	if err != nil {
		t.Fatalf("failed to list web application firewall policy: %+v", err)
	}
	t.Logf("listed web application firewall policy")

	err = ListAllWebApplicationFirewallPolicy(ctx)
	if err != nil {
		t.Fatalf("failed to list all web application firewall policy: %+v", err)
	}
	t.Logf("listed all web application firewall policy")

	err = DeleteWebApplicationFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to delete web application firewall policy: %+v", err)
	}
	t.Logf("deleted web application firewall policy")
}
