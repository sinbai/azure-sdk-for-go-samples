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

func getRouteFilterRuleClient() armnetwork.RouteFilterRulesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewRouteFilterRulesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

//Creates or updates a route in the specified route filter.
func CreateRouteFilterRule(ctx context.Context, routeFilterName string, routeName string) error {
	client := getRouteFilterRuleClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		routeFilterName,
		routeName,
		armnetwork.RouteFilterRule{
			Properties: &armnetwork.RouteFilterRulePropertiesFormat{
				Access:              armnetwork.AccessAllow.ToPtr(),
				Communities:         &[]*string{to.StringPtr("12076:51004")},
				RouteFilterRuleType: armnetwork.RouteFilterRuleTypeCommunity.ToPtr(),
			},
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

// Gets the specified route filter rule.
func GetRouteFilterRule(ctx context.Context, routeFilterName string, ruleName string) error {
	client := getRouteFilterRuleClient()
	_, err := client.Get(ctx, config.GroupName(), routeFilterName, ruleName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all RouteFilterRules in a route filter.
func ListRouteFilterByRouteFilter(ctx context.Context, routeFilterName string) error {
	client := getRouteFilterRuleClient()
	pager := client.ListByRouteFilter(config.GroupName(), routeFilterName, nil)
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

// Deletes the specified rule from a route filter.
func DeleteRouteFilterRule(ctx context.Context, routeFilterName string, ruleName string) error {
	client := getRouteFilterRuleClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), routeFilterName, ruleName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
