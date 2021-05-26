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

	publicIPAddressParameters := armnetwork.PublicIPAddress{
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
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	loadBalancerUrl := "/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName
	loadBalancerParameters := armnetwork.LoadBalancer{
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
							ID: to.StringPtr(loadBalancerUrl + "/backendAddressPools/" + backendAddressPoolName),
						},
						BackendPort:         to.Int32Ptr(80),
						DisableOutboundSnat: to.BoolPtr(true),
						EnableFloatingIP:    to.BoolPtr(true),
						EnableTCPReset:      new(bool),
						FrontendIPConfiguration: &armnetwork.SubResource{
							ID: to.StringPtr(loadBalancerUrl + "/frontendIPConfigurations/" + frontendIpConfigurationName),
						},
						FrontendPort:         to.Int32Ptr(80),
						IdleTimeoutInMinutes: to.Int32Ptr(15),
						LoadDistribution:     armnetwork.LoadDistributionDefault.ToPtr(),
						Probe: &armnetwork.SubResource{
							ID: to.StringPtr(loadBalancerUrl + "/probes/" + probeName),
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
							ID: to.StringPtr(loadBalancerUrl + "/backendAddressPools/" + backendAddressPoolName),
						},
						FrontendIPConfigurations: &[]*armnetwork.SubResource{
							{
								ID: to.StringPtr(loadBalancerUrl + "/frontendIPConfigurations/" + frontendIpConfigurationName),
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

	loadBalancerId, err := CreateLoadBalancer(ctx, loadBalancerName, loadBalancerParameters)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}
	t.Logf("created load balancer")

	inboundNatRuleParameters := armnetwork.InboundNatRule{
		Properties: &armnetwork.InboundNatRulePropertiesFormat{
			BackendPort:      to.Int32Ptr(3389),
			EnableFloatingIP: to.BoolPtr(false),
			EnableTCPReset:   to.BoolPtr(false),
			FrontendIPConfiguration: &armnetwork.SubResource{
				ID: to.StringPtr(loadBalancerId + "/frontendIPConfigurations/" + frontendIpConfigurationName),
			},
			FrontendPort:         to.Int32Ptr(3390),
			IdleTimeoutInMinutes: to.Int32Ptr(4),
			Protocol:             armnetwork.TransportProtocolTCP.ToPtr(),
		},
	}
	err = CreateInboundNatRule(ctx, loadBalancerName, inboundNatRuleName, inboundNatRuleParameters)
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
