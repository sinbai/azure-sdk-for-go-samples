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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestLoadBalancer(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	inboundNatRuleName := config.AppendRandomSuffix("inboundnatrule")
	loadBalancingRuleName := config.AppendRandomSuffix("loadbalancingrule")
	outBoundRuleName := config.AppendRandomSuffix("outboundrule")
	probeName := config.AppendRandomSuffix("probe")
	frontendIpConfigurationName := config.AppendRandomSuffix("frontendipconfiguration")
	backendAddressPoolName := config.AppendRandomSuffix("backendaddresspool")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	publicIPAddress := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			IdleTimeoutInMinutes:     to.Int32Ptr(10),
			PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
		},
		SKU: &armnetwork.PublicIPAddressSKU{
			Name: armnetwork.PublicIPAddressSKUNameStandard.ToPtr(),
		},
	}
	err = CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddress)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	err = CreateLoadBalancer(ctx, loadBalancerName, publicIpAddressName, frontendIpConfigurationName, backendAddressPoolName,
		probeName, loadBalancingRuleName, outBoundRuleName)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}
	t.Logf("created load balancer")

	err = CreateInboundNatRule(ctx, loadBalancerName, inboundNatRuleName, frontendIpConfigurationName)
	if err != nil {
		t.Fatalf("failed to get load balancer inbound nat rule: %+v", err)
	}

	err = GetLoadBalancerFrontendIPConfiguration(ctx, loadBalancerName, frontendIpConfigurationName)
	if err != nil {
		t.Fatalf("failed to get load balancer frontend IP configuration: %+v", err)
	}

	err = GetLoadBalancerBackendAddressPool(ctx, loadBalancerName, backendAddressPoolName)
	if err != nil {
		t.Fatalf("failed to get load balancer backend address pool: %+v", err)
	}

	err = GetLoadBalancerLoadBalancingRule(ctx, loadBalancerName, loadBalancingRuleName)
	if err != nil {
		t.Fatalf("failed to get the specified load balancer load balancing rule: %+v", err)
	}
	t.Logf("got the specified load balancer load balancing rule")

	err = GetLoadBalancerOutboundRule(ctx, loadBalancerName, outBoundRuleName)
	if err != nil {
		t.Fatalf("failed to get the specified load balancer outbound rule.: %+v", err)
	}
	t.Logf("got the specified load balancer outbound rule.")

	err = ListLoadBalancerFrontendIPConfiguration(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to get all the load balancer frontend IP configurations: %+v", err)
	}
	t.Logf("listed all the load balancer frontend IP configurations.")

	err = GetLoadBalancerProbe(ctx, loadBalancerName, probeName)
	if err != nil {
		t.Fatalf("failed to get load balancer probe: %+v", err)
	}
	t.Logf("got load balancer probe")

	err = ListLoadBalancerBackendAddressPool(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to get all the load balancer backed address pools: %+v", err)
	}
	t.Logf("listed all the load balancer backed address pools")

	err = ListLoadBalancerLoadBalancingRule(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to list all the load balancing rules in a load balancer: %+v", err)
	}
	t.Logf("listed all the load balancing rules in a load balancer")

	err = ListLoadBalancerNetworkInterface(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to list associated load balancer network interfaces: %+v", err)
	}
	t.Logf("listed associated load balancer network interfaces")

	err = ListLoadBalancerNetworkInterface(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to list associated load balancer network interfaces: %+v", err)
	}
	t.Logf("listed associated load balancer network interfaces")

	err = ListLoadBalancerOutboundRule(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to list all the outbound rules in a load balancer: %+v", err)
	}
	t.Logf("listed all the outbound rules in a load balancer")

	err = ListLoadBalancerProbe(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to list load balancer probe: %+v", err)
	}
	t.Logf("listed load balancer probe")

	err = GetLoadBalancer(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to get the specified load balancer in a specified resource group: %+v", err)
	}
	t.Logf("got the specified load balancer in a specified resource group")

	err = ListLoadBalancer(ctx)
	if err != nil {
		t.Fatalf("failed to list all the load balancers in a resource group: %+v", err)
	}
	t.Logf("listed all the load balancers in a resource group")

	err = ListAllLoadBalancer(ctx)
	if err != nil {
		t.Fatalf("failed to list all the load balancer in a subscription: %+v", err)
	}
	t.Logf("listed all the load balancer in a subscription")

	err = UpdateLoadBalancerTags(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to update tags for load balancer: %+v", err)
	}
	t.Logf("updated load balancer tags")

	err = DeleteLoadBalancer(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to delete load balancer: %+v", err)
	}
	t.Logf("deleted load balancer")
}
