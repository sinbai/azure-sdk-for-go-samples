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

func TestExpressRouteCircuitPeering(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	expressRouteCircuitPeeringName := "AzurePrivatePeering"
	expressRouteCircuitName := config.AppendRandomSuffix("expressroutecircuit")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateExpressRouteCircuit(ctx, expressRouteCircuitName)
	if err != nil {
		t.Fatalf("failed to create express route circuit: % +v", err)
	}

	err = CreateExpressRouteCircuitPeering(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName)
	if err != nil {
		t.Fatalf("failed to create express route circuit peering: % +v", err)
	}
	t.Logf("created express route circuit peering")

	err = GetExpressRouteCircuitPeering(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName)
	if err != nil {
		t.Fatalf("failed to get express route circuit peering: %+v", err)
	}
	t.Logf("got express route circuit peering")

	err = ListExpressRouteCircuitPeering(ctx, expressRouteCircuitName)
	if err != nil {
		t.Fatalf("failed to list express route circuit peering: %+v", err)
	}
	t.Logf("listed express route circuit peering")

	err = DeleteExpressRouteCircuitPeering(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName)
	if err != nil {
		t.Fatalf("failed to delete express route circuit peering: %+v", err)
	}
	t.Logf("deleted express route circuit peering")
}
