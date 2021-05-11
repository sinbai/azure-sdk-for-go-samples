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
)

func TestFirewallPolicyRullCollectionGroup(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	firewallPolicyName := config.AppendRandomSuffix("firewallpolicy")
	firewallPolicyRuleCollectionGroupName := config.AppendRandomSuffix("firewallpolicyRuleCollectionGroup")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateFirewallPolicy(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to create firewall policy: % +v", err)
	}
	t.Logf("created firewall policy")

	body := `{
		"priority": 100,
		"ruleCollections": [
		{
			"ruleCollectionType": "FirewallPolicyNatRuleCollection",
			"priority": 100,
			"name": "Example-Nat-Rule-Collection",
			"action": {
			"type": "DNAT"
			},
			"rules": [
			{
				"ruleType": "NatRule",
				"name": "nat-rule1",
				"translatedFqdn": "internalhttp.server.net",
				"translatedPort": "8080",
				"ipProtocols": [
				"TCP",
				"UDP"
				],
				"sourceAddresses": [
				"2.2.2.2"
				],
				"sourceIpGroups": [],
				"destinationAddresses": [
				"152.23.32.23"
				],
				"destinationPorts": [
				"8080"
				]
			}
			]
		}
		]
		}`
	err = CreateFirewallPolicyRuleCollectionGroup(ctx, firewallPolicyName, firewallPolicyRuleCollectionGroupName, body)
	if err != nil {
		t.Fatalf("failed to create specified firewall policy rule collection group: % +v", err)
	}
	t.Logf("created specified firewall policy rule collection group")

	err = GetFirewallPolicyRuleCollectionGroup(ctx, firewallPolicyName, firewallPolicyRuleCollectionGroupName)
	if err != nil {
		t.Fatalf("failed to get firewall policy rule collection group: %+v", err)
	}
	t.Logf("got firewall policy rule collection group")

	err = ListFirewallPolicyRuleCollectionGroup(ctx, firewallPolicyName)
	if err != nil {
		t.Fatalf("failed to list firewall policy rule collection group: %+v", err)
	}
	t.Logf("listed firewall policy rule collection group")

	err = DeleteFirewallPolicyRuleCollectionGroup(ctx, firewallPolicyName, firewallPolicyRuleCollectionGroupName)
	if err != nil {
		t.Fatalf("failed to delete firewall policy rule collection group: %+v", err)
	}
	t.Logf("deleted firewall policy rule collection group")
}
