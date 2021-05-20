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

func TestPrivateEndpoint(t *testing.T) {
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

	virtualNetworkPro := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkPro)
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
	subnet2ID, err := CreateSubnet(ctx, virtualNetworkName, subNetName2, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	loadBalancerPro := armnetwork.LoadBalancer{
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

	loadBalancerId, err := CreateLoadBalancer(ctx, loadBalancerName, loadBalancerPro)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}

	privateLinkServicePro := armnetwork.PrivateLinkService{
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
	privateLinkServiceId, err := CreatePrivateLinkService(ctx, privateLinkServiceName, privateLinkServicePro)
	if err != nil {
		t.Fatalf("failed to create private link service: % +v", err)
	}

	privateEndpointPro := armnetwork.PrivateEndpoint{
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

	err = CreatePrivateEndpoint(ctx, privateEndpointName, privateEndpointPro)
	if err != nil {
		t.Fatalf("failed to create private endpoint: % +v", err)
	}
	t.Logf("created private endpoint")

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
