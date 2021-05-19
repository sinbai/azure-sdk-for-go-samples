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

	publicIPAddressPro := armnetwork.PublicIPAddress{
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
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressPro)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	loadBalancerPro := armnetwork.LoadBalancer{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			BackendAddressPools: &[]*armnetwork.BackendAddressPool{
				{
					Name: &backendAddressPoolName,
				},
			},
			FrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{
				{
					Name: &frontendIpConfigurationName,
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							Resource: armnetwork.Resource{
								ID: &publicIpAddressId,
							},
						},
					},
				},
			},
			LoadBalancingRules: &[]*armnetwork.LoadBalancingRule{
				{
					Name: &loadBalancingRuleName,
					Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
						BackendAddressPool: &armnetwork.SubResource{
							ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName + "/backendAddressPools/" + backendAddressPoolName),
						},
						BackendPort:         to.Int32Ptr(80),
						DisableOutboundSnat: to.BoolPtr(true),
						EnableFloatingIP:    to.BoolPtr(true),
						EnableTCPReset:      new(bool),
						FrontendIPConfiguration: &armnetwork.SubResource{
							ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName + "/frontendIPConfigurations/" + frontendIpConfigurationName),
						},
						FrontendPort:         to.Int32Ptr(80),
						IdleTimeoutInMinutes: to.Int32Ptr(15),
						LoadDistribution:     armnetwork.LoadDistributionDefault.ToPtr(),
						Probe: &armnetwork.SubResource{
							ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName + "/probes/" + probeName),
						},
						Protocol: armnetwork.TransportProtocolTCP.ToPtr(),
					},
				},
			},
			OutboundRules: &[]*armnetwork.OutboundRule{
				{
					Name: &outBoundRuleName,
					Properties: &armnetwork.OutboundRulePropertiesFormat{
						BackendAddressPool: &armnetwork.SubResource{
							ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName + "/backendAddressPools/" + backendAddressPoolName),
						},
						FrontendIPConfigurations: &[]*armnetwork.SubResource{
							{
								ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName + "/frontendIPConfigurations/" + frontendIpConfigurationName),
							},
						},
						Protocol: armnetwork.LoadBalancerOutboundRuleProtocolAll.ToPtr(),
					},
				},
			},
			Probes: &[]*armnetwork.Probe{
				{
					Name: &probeName,
					Properties: &armnetwork.ProbePropertiesFormat{
						IntervalInSeconds: to.Int32Ptr(15),
						NumberOfProbes:    to.Int32Ptr(2),
						Port:              to.Int32Ptr(80),
						Protocol:          armnetwork.ProbeProtocolHTTP.ToPtr(),
						RequestPath:       to.StringPtr("healthcheck.aspx"),
					},
				},
			},
		},
		SKU: &armnetwork.LoadBalancerSKU{
			Name: armnetwork.LoadBalancerSKUNameStandard.ToPtr(),
		},
	}

	err = CreateLoadBalancer(ctx, loadBalancerName, loadBalancerPro)
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
