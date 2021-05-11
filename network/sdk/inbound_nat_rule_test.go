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

func TestInboundNatRule(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	inboundNatRuleName := config.AppendRandomSuffix("inboundnatrule")
	probeName := config.AppendRandomSuffix("probe")
	frontendIpConfigurationName := config.AppendRandomSuffix("frontendipconfiguration")
	backendAddressPoolName := config.AppendRandomSuffix("backendaddresspool")
	loadBalancingRuleName := config.AppendRandomSuffix("loadbalancingrule")
	outBoundRuleName := config.AppendRandomSuffix("outboundrule")

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

	err = CreateLoadBalancer(ctx, loadBalancerName, publicIpAddressName, frontendIpConfigurationName,
		backendAddressPoolName, probeName, loadBalancingRuleName, outBoundRuleName)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}
	t.Logf("created load balancer")

	err = CreateInboundNatRule(ctx, loadBalancerName, inboundNatRuleName, frontendIpConfigurationName)
	if err != nil {
		t.Fatalf("failed to get load balancer inbound nat rule: %+v", err)
	}

	err = GetInboundNatRule(ctx, loadBalancerName, inboundNatRuleName)
	if err != nil {
		t.Fatalf("failed to list the specified load balancer inbound nat rule: %+v", err)
	}
	t.Logf("listed the specified load balancer inbound nat rule")

	err = ListInboundNatRule(ctx, loadBalancerName)
	if err != nil {
		t.Fatalf("failed to list all the inbound nat rules in a load balancer: %+v", err)
	}
	t.Logf("listed all the inbound nat rules in a load balancer")

	err = DeleteInboundNatRule(ctx, loadBalancerName, inboundNatRuleName)
	if err != nil {
		t.Fatalf("failed to delete the specified load balancer inbound nat rule: %+v", err)
	}
	t.Logf("deleted the specified load balancer inbound nat rule")
}
