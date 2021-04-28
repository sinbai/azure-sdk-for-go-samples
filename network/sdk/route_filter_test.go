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

func TestRouteFilter(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	routeFilterName := config.AppendRandomSuffix("routefilter")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateRouteFilter(ctx, routeFilterName)
	if err != nil {
		t.Fatalf("failed to create route filter: %+v", err)
	}
	t.Logf("created route filter")

	err = GetRouteFilter(ctx, routeFilterName)
	if err != nil {
		t.Fatalf("failed to get route filter: %+v", err)
	}
	t.Logf("got route filter")

	err = ListRouteFilter(ctx)
	if err != nil {
		t.Fatalf("failed to list route filter: %+v", err)
	}
	t.Logf("listed route filter")

	err = ListRouteFilterByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list route filters by resource group: %+v", err)
	}
	t.Logf("listed route filter by resource group")

	err = UpdateRouteFilterTags(ctx, routeFilterName)
	if err != nil {
		t.Fatalf("failed to update tags for route filter: %+v", err)
	}
	t.Logf("updated route filter tags")

	err = DeleteRouteFilter(ctx, routeFilterName)
	if err != nil {
		t.Fatalf("failed to delete route filter: %+v", err)
	}
	t.Logf("deleted route filter")

}
