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

func TestPrivateLinkService(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	privateLinkServiceName := config.AppendRandomSuffix("privatelinkservice")
	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName1 := config.AppendRandomSuffix("subnet1")
	subNetName2 := config.AppendRandomSuffix("subnet2")
	ipConfigName := config.AppendRandomSuffix("ipconfig")
	privateEndpointName := config.AppendRandomSuffix("privateendpoint")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateVirtualNetwork(ctx, virtualNetworkName)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	body := `{
	"addressPrefix": "10.0.1.0/24",
	"privateLinkServiceNetworkPolicies": "Disabled"
	}`
	subnet1ID, err := CreateSubnet(ctx, virtualNetworkName, subNetName1, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	body = `{
		"addressPrefix": "10.0.0.0/24",
		"privateEndpointNetworkPolicies": "Disabled"
		}`
	_, err = CreateSubnet(ctx, virtualNetworkName, subNetName2, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	err = createLoadBalancer(ctx, loadBalancerName, ipConfigName, subnet1ID)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}

	err = CreatePrivateLinkService(ctx, privateLinkServiceName, virtualNetworkName, loadBalancerName, ipConfigName, subNetName1)
	if err != nil {
		t.Fatalf("failed to create private link service: % +v", err)
	}
	t.Logf("created private like service")

	err = CreatePrivateEndpoint(ctx, privateEndpointName, privateLinkServiceName, virtualNetworkName, subNetName2)
	if err != nil {
		t.Fatalf("failed to create private endpoint: % +v", err)
	}

	peConnectionName, err := GetPrivateLinkService(ctx, privateLinkServiceName)
	if err != nil {
		t.Fatalf("failed to get private link service: %+v", err)
	}
	t.Logf("got private link service")

	err = UpdatePrivateEndpointConnection(ctx, privateLinkServiceName, peConnectionName)
	if err != nil {
		t.Fatalf("failed to update tags for private link service: %+v", err)
	}
	t.Logf("updated private link service tags")

	err = GetPrivateEndpointConnection(ctx, privateLinkServiceName, peConnectionName)
	if err != nil {
		t.Fatalf("failed to get private endpoint connection: %+v", err)
	}
	t.Logf("got private endpoint connection")

	err = ListPrivateEndpointConnections(ctx, privateLinkServiceName)
	if err != nil {
		t.Fatalf("failed to list private endpoint connection: %+v", err)
	}
	t.Logf("listed private endpoint connection")

	err = ListAutoApprovedPrivateLinkServicesByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list approved private link services by resource group: %+v", err)
	}
	t.Logf("listed approved private link services by resource group")

	_, err = GetPrivateLinkService(ctx, privateLinkServiceName)
	if err != nil {
		t.Fatalf("failed to get private link service: %+v", err)
	}
	t.Logf("got private link service")

	err = ListPrivateEndpointConnections(ctx, privateLinkServiceName)
	if err != nil {
		t.Fatalf("failed to list private endpoint connection: %+v", err)
	}
	t.Logf("listed private endpoint connection")

	err = ListBySubscription(ctx)
	if err != nil {
		t.Fatalf("failed to list private link services by subscription: %+v", err)
	}
	t.Logf("listed private link services by subscription")

	err = BeginDeletePrivateEndpointConnection(ctx, privateLinkServiceName, peConnectionName)
	if err != nil {
		t.Fatalf("failed to begin delete private endpoint connection: %+v", err)
	}
	t.Logf("begin deleted private endpoint connection")

	err = DeletePrivateLinkService(ctx, privateLinkServiceName)
	if err != nil {
		t.Fatalf("failed to delete private link service: %+v", err)
	}
	t.Logf("deleted private link service")
}
