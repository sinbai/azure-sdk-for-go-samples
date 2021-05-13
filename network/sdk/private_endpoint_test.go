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

func TestPrivateEndpoint(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	privateLinkServiceName := config.AppendRandomSuffix("privatelinkservice")
	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName1 := config.AppendRandomSuffix("subnet1")
	subNetName2 := config.AppendRandomSuffix("subnet2")
	ipConfigName := config.AppendRandomSuffix("ipconfig")
	privateEndpointName := config.AppendRandomSuffix("privateendpoint")
	privateZoneName := config.AppendRandomSuffix("www.zone1.com")

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

	err = CreatePrivateEndpoint(ctx, privateEndpointName, privateLinkServiceName, virtualNetworkName, subNetName2)
	if err != nil {
		t.Fatalf("failed to create private endpoint: % +v", err)
	}
	t.Logf("created private endpoint")

	_, err = CreatePrivateZone(ctx, privateZoneName)
	if err != nil {
		t.Fatalf("failed to create private zone: % +v", err)
	}
	t.Logf("created private zone")

	err = GetPrivateEndpoint(ctx, privateEndpointName)
	if err != nil {
		t.Fatalf("failed to get private endpoint: %+v", err)
	}
	t.Logf("got private endpoint")

	err = ListPrivateEndpoint(ctx)
	if err != nil {
		t.Fatalf("failed to list private endpoint: %+v", err)
	}
	t.Logf("listed private endpoint")

	err = ListAllPrivateEndpointBySubscription(ctx)
	if err != nil {
		t.Fatalf("failed to list all private endpoint: %+v", err)
	}
	t.Logf("listed all private endpoint")

	err = DeletePrivateEndpoint(ctx, privateEndpointName)
	if err != nil {
		t.Fatalf("failed to delete private endpoint: %+v", err)
	}
	t.Logf("deleted private endpoint")

}
