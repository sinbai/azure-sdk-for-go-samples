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

func TestPrivateDnsZoneGroup(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	privateDnsZoneGroupName := config.AppendRandomSuffix("privatednszonegroup")

	privateLinkServiceName := config.AppendRandomSuffix("privatelinkservice")
	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName1 := config.AppendRandomSuffix("subnet")
	subNetName2 := config.AppendRandomSuffix("subnet")
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

	err = CreateVirtualNetwork(ctx, virtualNetworkName, "10.0.0.0/16")
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

	err = CreatePrivateDnsZoneGroup(ctx, privateEndpointName, privateDnsZoneGroupName, privateZoneName)
	if err != nil {
		t.Fatalf("failed to create private dns zone group: % +v", err)
	}
	t.Logf("created private dns zone group")

	err = GetPrivateDnsZoneGroup(ctx, privateEndpointName, privateDnsZoneGroupName)
	if err != nil {
		t.Fatalf("failed to get private dns zone group: %+v", err)
	}
	t.Logf("got private dns zone group")

	err = ListPrivateDnsZoneGroup(ctx, privateEndpointName)
	if err != nil {
		t.Fatalf("failed to list private dns zone group: %+v", err)
	}
	t.Logf("listed private dns zone group")

	err = DeletePrivateDnsZoneGroup(ctx, privateEndpointName, privateDnsZoneGroupName)
	if err != nil {
		t.Fatalf("failed to delete private dns zone group: %+v", err)
	}
	t.Logf("deleted private dns zone group")

}
