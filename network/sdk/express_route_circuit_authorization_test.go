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

func TestExpressRouteCircuitAuthorization(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	expressRouteCircuitAuthorizationName := config.AppendRandomSuffix("expressroutecircuitauthorization")
	expressRouteCircuitName := config.AppendRandomSuffix("expressroutecircuit")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	expressRouteCircuitParameters := armnetwork.ExpressRouteCircuit{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.ExpressRouteCircuitPropertiesFormat{
			ServiceProviderProperties: &armnetwork.ExpressRouteCircuitServiceProviderProperties{
				BandwidthInMbps:     to.Int32Ptr(200),
				PeeringLocation:     to.StringPtr("Silicon Valley Test"),
				ServiceProviderName: to.StringPtr("Equinix Test"),
			},
		},
		SKU: &armnetwork.ExpressRouteCircuitSKU{
			Family: armnetwork.ExpressRouteCircuitSKUFamilyMeteredData.ToPtr(),
			Name:   to.StringPtr("Standard_MeteredData"),
			Tier:   armnetwork.ExpressRouteCircuitSKUTierStandard.ToPtr(),
		},
	}
	_, err = CreateExpressRouteCircuit(ctx, expressRouteCircuitName, expressRouteCircuitParameters)
	if err != nil {
		t.Fatalf("failed to create express route circuit: % +v", err)
	}

	err = CreateExpressRouteCircuitAuthorization(ctx, expressRouteCircuitName, expressRouteCircuitAuthorizationName)
	if err != nil {
		t.Fatalf("failed to create express route circuit authorization: % +v", err)
	}
	t.Logf("created express route circuit authorization")

	err = GetExpressRouteCircuitAuthorization(ctx, expressRouteCircuitName, expressRouteCircuitAuthorizationName)
	if err != nil {
		t.Fatalf("failed to get express route circuit authorization: %+v", err)
	}
	t.Logf("got express route circuit authorization")

	err = ListExpressRouteCircuitAuthorization(ctx, expressRouteCircuitName)
	if err != nil {
		t.Fatalf("failed to list express route circuit authorization: %+v", err)
	}
	t.Logf("listed express route circuit authorization")

	err = DeleteExpressRouteCircuitAuthorization(ctx, expressRouteCircuitName, expressRouteCircuitAuthorizationName)
	if err != nil {
		t.Fatalf("failed to delete express route circuit authorization: %+v", err)
	}
	t.Logf("deleted express route circuit authorization")
}
