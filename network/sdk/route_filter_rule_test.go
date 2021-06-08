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

func TestRouteFilterRule(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	routeFilterName := config.AppendRandomSuffix("routefilter")
	ruleName := config.AppendRandomSuffix("rule")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	routeFilterParameters := armnetwork.RouteFilter{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
	}
	err = CreateRouteFilter(ctx, routeFilterName, routeFilterParameters)
	if err != nil {
		t.Fatalf("failed to create route filter: %+v", err)
	}

	routeFilterRuleParameters := armnetwork.RouteFilterRule{
		Properties: &armnetwork.RouteFilterRulePropertiesFormat{
			Access:              armnetwork.AccessAllow.ToPtr(),
			Communities:         &[]*string{to.StringPtr("12076:51004")},
			RouteFilterRuleType: armnetwork.RouteFilterRuleTypeCommunity.ToPtr(),
		},
	}
	err = CreateRouteFilterRule(ctx, routeFilterName, ruleName, routeFilterRuleParameters)
	if err != nil {
		t.Fatalf("failed to create route filter rule: %+v", err)
	}
	t.Logf("created route filter rule")

	err = GetRouteFilterRule(ctx, routeFilterName, ruleName)
	if err != nil {
		t.Fatalf("failed to get route filter rule: %+v", err)
	}
	t.Logf("got route filter rule")

	err = ListRouteFilterByRouteFilter(ctx, routeFilterName)
	if err != nil {
		t.Fatalf("failed to list rule by route filter: %+v", err)
	}
	t.Logf("listed rule by route filter")

	err = DeleteRouteFilterRule(ctx, routeFilterName, ruleName)
	if err != nil {
		t.Fatalf("failed to delete route filter rule: %+v", err)
	}
	t.Logf("deleted route filter rule")

}
