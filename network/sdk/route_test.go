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

func TestRoute(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	routeTableName := config.AppendRandomSuffix("routetable")
	routeName := config.AppendRandomSuffix("route")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateRouteTable(ctx, routeTableName)
	if err != nil {
		t.Fatalf("failed to create route table: % +v", err)
	}
	t.Logf("created route table")

	routeParameters := armnetwork.Route{
		Properties: &armnetwork.RoutePropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.3.0/24"),
			NextHopType:   armnetwork.RouteNextHopTypeVirtualNetworkGateway.ToPtr(),
		},
	}
	err = CreateRoute(ctx, routeTableName, routeName, routeParameters)
	if err != nil {
		t.Fatalf("failed to create route: % +v", err)
	}
	t.Logf("created route")

	err = GetRoute(ctx, routeTableName, routeName)
	if err != nil {
		t.Fatalf("failed to get route: %+v", err)
	}
	t.Logf("got route")

	err = ListRoute(ctx, routeTableName)
	if err != nil {
		t.Fatalf("failed to list route: %+v", err)
	}
	t.Logf("listed route")

	err = DeleteRoute(ctx, routeTableName, routeName)
	if err != nil {
		t.Fatalf("failed to delete route: %+v", err)
	}
	t.Logf("deleted route")
}
