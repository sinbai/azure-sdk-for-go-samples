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

func TestRouteTable(t *testing.T) {
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

	err = CreateRouteTable(ctx, routeTableName, routeName)
	if err != nil {
		t.Fatalf("failed to create route table: % +v", err)
	}
	t.Logf("created route table")

	err = CreateRoute(ctx, routeTableName, routeName)
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

	err = GetRouteTable(ctx, routeTableName)
	if err != nil {
		t.Fatalf("failed to get route table: %+v", err)
	}
	t.Logf("got route table")
	err = ListRouteTable(ctx)
	if err != nil {
		t.Fatalf("failed to list route table: %+v", err)
	}
	t.Logf("listed route table")
	err = ListAllRouteTable(ctx)
	if err != nil {
		t.Fatalf("failed to list all route table: %+v", err)
	}
	t.Logf("listed all route table")
	err = UpdateRouteTableTags(ctx, routeTableName)
	if err != nil {
		t.Fatalf("failed to update tags for route table: %+v", err)
	}
	t.Logf("updated route table tags")

	err = DeleteRoute(ctx, routeTableName, routeName)
	if err != nil {
		t.Fatalf("failed to delete route: %+v", err)
	}
	t.Logf("deleted route")

	err = DeleteRouteTable(ctx, routeTableName)
	if err != nil {
		t.Fatalf("failed to delete route table: %+v", err)
	}
	t.Logf("deleted route table")
}