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
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
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
	ipConfigurationName := config.AppendRandomSuffix("ipconfig")
	privateEndpointName := config.AppendRandomSuffix("privateendpoint")
	privateZoneName := config.AppendRandomSuffix("www.zone1.com")

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
			PrivateLinkServiceNetworkPolicies: to.StringPtr("Disable"),
		},
	}
	subnet1ID, err := CreateSubnet(ctx, virtualNetworkName, subNetName1, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	subnetParameters = armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                     to.StringPtr("10.0.0.0/24"),
			PrivateLinkServiceNetworkPolicies: to.StringPtr("Disable"),
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

	err = CreatePrivateEndpoint(ctx, privateEndpointName, privateEndpointParameters)
	if err != nil {
		t.Fatalf("failed to create private endpoint: % +v", err)
	}
	t.Logf("created private endpoint")

	privateZoneParameters := privatedns.PrivateZone{
		Location: to.StringPtr("global"),
	}
	privateZone, err := CreatePrivateZone(ctx, privateZoneName, privateZoneParameters)
	if err != nil {
		t.Fatalf("failed to create private zone: % +v", err)
	}
	t.Logf("created private zone")

	privateDNSZoneGroupParameters := armnetwork.PrivateDNSZoneGroup{
		Name: &privateDnsZoneGroupName,
		Properties: &armnetwork.PrivateDNSZoneGroupPropertiesFormat{
			PrivateDNSZoneConfigs: &[]*armnetwork.PrivateDNSZoneConfig{{
				Name: &privateZoneName,
				Properties: &armnetwork.PrivateDNSZonePropertiesFormat{
					PrivateDNSZoneID: privateZone.ID,
				},
			}},
		},
	}
	err = CreatePrivateDnsZoneGroup(ctx, privateEndpointName, privateDnsZoneGroupName, privateDNSZoneGroupParameters)
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
