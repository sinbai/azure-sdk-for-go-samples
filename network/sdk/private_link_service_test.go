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

func TestPrivateLinkService(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	privateLinkServiceName := config.AppendRandomSuffix("privatelinkservice")
	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName1 := config.AppendRandomSuffix("subnet1")
	subNetName2 := config.AppendRandomSuffix("subnet2")
	ipConfigurationName := config.AppendRandomSuffix("ipconfig")
	privateEndpointName := config.AppendRandomSuffix("privateendpoint")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualNetworkParameters := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                     to.StringPtr("10.0.1.0/24"),
			PrivateLinkServiceNetworkPolicies: to.StringPtr("Disabled"),
		},
	}
	subnet1ID, err := CreateSubnet(ctx, virtualNetworkName, subNetName1, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	subnetParameters = armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                  to.StringPtr("10.0.0.0/24"),
			PrivateEndpointNetworkPolicies: to.StringPtr("Disabled"),
		},
	}
	subnet2ID, err := CreateSubnet(ctx, virtualNetworkName, subNetName2, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	loadBalancerParameters := armnetwork.LoadBalancer{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{
				{
					Name: &ipConfigurationName,
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.Subnet{
							SubResource: armnetwork.SubResource{
								ID: &subnet1ID,
							},
						},
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

	privateLinkServiceParameters := armnetwork.PrivateLinkService{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.PrivateLinkServiceProperties{
			AutoApproval: &armnetwork.PrivateLinkServicePropertiesAutoApproval{
				ResourceSet: armnetwork.ResourceSet{
					Subscriptions: &[]*string{to.StringPtr(config.SubscriptionID())},
				},
			},
			Fqdns: &[]*string{to.StringPtr("fqdn1"),
				to.StringPtr("fqdn2"),
				to.StringPtr("fqdn3")},
			IPConfigurations: &[]*armnetwork.PrivateLinkServiceIPConfiguration{{
				Name: &ipConfigurationName,
				Properties: &armnetwork.PrivateLinkServiceIPConfigurationProperties{
					PrivateIPAddress:          to.StringPtr("10.0.1.5"),
					PrivateIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
					PrivateIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
					Subnet: &armnetwork.Subnet{
						SubResource: armnetwork.SubResource{
							ID: &subnet1ID,
						},
					},
				},
			}},
			LoadBalancerFrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{{
				SubResource: armnetwork.SubResource{
					ID: to.StringPtr(loadBalancerId + "/frontendIPConfigurations/" + ipConfigurationName),
				},
			}},
			Visibility: &armnetwork.PrivateLinkServicePropertiesVisibility{
				ResourceSet: armnetwork.ResourceSet{
					Subscriptions: &[]*string{to.StringPtr(config.SubscriptionID())},
				},
			},
		},
	}

	privateLinkServiceId, err := CreatePrivateLinkService(ctx, privateLinkServiceName, privateLinkServiceParameters)
	if err != nil {
		t.Fatalf("failed to create private link service: % +v", err)
	}
	t.Logf("created private like service")

	privateEndpointParameters := armnetwork.PrivateEndpoint{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.PrivateEndpointProperties{
			PrivateLinkServiceConnections: &[]*armnetwork.PrivateLinkServiceConnection{{
				Name: &privateLinkServiceName,
				Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
					PrivateLinkServiceID: &privateLinkServiceId,
				},
			}},
			Subnet: &armnetwork.Subnet{
				SubResource: armnetwork.SubResource{
					ID: &subnet2ID,
				},
			},
		},
	}

	_, err = CreatePrivateEndpoint(ctx, privateEndpointName, privateEndpointParameters)
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
