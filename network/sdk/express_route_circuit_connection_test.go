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

func TestExpressRouteCircuitConnection(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	expressRouteCircuitConnectionName := config.AppendRandomSuffix("expressroutecircuitconnection")
	expressRouteCircuitName := config.AppendRandomSuffix("expressroutecircuit")
	expressRouteCircuitName2 := config.AppendRandomSuffix("expressroutecircuit2")
	expressRouteCircuitPeeringName := "AzurePrivatePeering"

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

	err = CreateExpressRouteCircuit(ctx, expressRouteCircuitName2)
	if err != nil {
		t.Fatalf("failed to create express route circuit2: % +v", err)
	}

	err = CreateExpressRouteCircuitPeering(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName)
	if err != nil {
		t.Fatalf("failed to create express route circuit peering: % +v", err)
	}

	err = CreateExpressRouteCircuitConnection(ctx, expressRouteCircuitName, expressRouteCircuitName2,
		expressRouteCircuitPeeringName, expressRouteCircuitConnectionName)
	if err != nil {
		t.Fatalf("failed to create express route circuit connection: % +v", err)
	}
	t.Logf("created express route circuit connection")

	err = GetExpressRouteCircuitConnection(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName, expressRouteCircuitConnectionName)
	if err != nil {
		t.Fatalf("failed to get express route circuit connection: %+v", err)
	}
	t.Logf("got express route circuit connection")

	err = ListExpressRouteCircuitConnection(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName)
	if err != nil {
		t.Fatalf("failed to list express route circuit connection: %+v", err)
	}
	t.Logf("listed express route circuit connection")

	err = DeleteExpressRouteCircuitConnection(ctx, expressRouteCircuitName, expressRouteCircuitPeeringName, expressRouteCircuitConnectionName)
	if err != nil {
		t.Fatalf("failed to delete express route circuit connection: %+v", err)
	}
	t.Logf("deleted express route circuit connection")
}
