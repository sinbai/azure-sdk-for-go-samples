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

func TestVirtualNetworkGatewayConnection(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualNetworkGatewayConnectionName := config.AppendRandomSuffix("virtualnetworkgatewayconnection")
	localNetworkGatewayName := config.AppendRandomSuffix("localvirtualnetworkgateway")
	virtualNetworkGatewayName := config.AppendRandomSuffix("virtualnetworkgateway")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	gatewaySubNetName := "GatewaySubnet"
	ipConfigName := config.AppendRandomSuffix("ipconfig")

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Name:     to.StringPtr(publicIpAddressName),
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
			IdleTimeoutInMinutes:     to.Int32Ptr(4),
		},
	}

	publicAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
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
			AddressPrefix: to.StringPtr("10.0.1.0/24"),
		},
	}
	subnetId, err := CreateSubnet(ctx, virtualNetworkName, gatewaySubNetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	virtualNetWorkGatewayParameters := armnetwork.VirtualNetworkGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
			Active: to.BoolPtr(false),
			BgpSettings: &armnetwork.BgpSettings{
				Asn:               to.Int64Ptr(65515),
				BgpPeeringAddress: to.StringPtr("10.0.1.30"),
				PeerWeight:        to.Int32Ptr(0),
			},
			CustomRoutes: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("101.168.0.6/32")},
			},
			EnableBgp:           to.BoolPtr(false),
			EnableDNSForwarding: to.BoolPtr(false),
			GatewayType:         armnetwork.VirtualNetworkGatewayTypeVPN.ToPtr(),
			IPConfigurations: &[]*armnetwork.VirtualNetworkGatewayIPConfiguration{{
				Name: &ipConfigName,
				Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
					PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
					PublicIPAddress: &armnetwork.SubResource{
						ID: &publicAddressId,
					},
					Subnet: &armnetwork.SubResource{
						ID: &subnetId,
					},
				},
			}},
			SKU: &armnetwork.VirtualNetworkGatewaySKU{
				Name: armnetwork.VirtualNetworkGatewaySKUNameVPNGw1.ToPtr(),
				Tier: armnetwork.VirtualNetworkGatewaySKUTierVPNGw1.ToPtr(),
			},
			VPNType: armnetwork.VPNTypeRouteBased.ToPtr(),
		},
	}

	gatewayId, err := CreateVirtualNetworkGateway(ctx, virtualNetworkGatewayName, virtualNetWorkGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network gateway: % +v", err)
	}
	t.Logf("created virtual network gateway")

	localNetworkGatewayParameters := armnetwork.LocalNetworkGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
			GatewayIPAddress: to.StringPtr("11.12.13.14"),
			LocalNetworkAddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.1.0.0/16")},
			},
		},
	}
	localGatewayId, err := CreateLocalNetworkGateway(ctx, localNetworkGatewayName, localNetworkGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create local network gateway: % +v", err)
	}

	gatewayConnectionParameters := armnetwork.VirtualNetworkGatewayConnection{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
			ConnectionProtocol: armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2.ToPtr(),
			ConnectionType:     armnetwork.VirtualNetworkGatewayConnectionTypeIPsec.ToPtr(),
			EnableBgp:          to.BoolPtr(false),
			LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
				Resource: armnetwork.Resource{
					ID: &localGatewayId,
				},
				Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
					GatewayIPAddress: to.StringPtr("10.1.0.1"),
					LocalNetworkAddressSpace: &armnetwork.AddressSpace{
						AddressPrefixes: &[]*string{to.StringPtr("10.1.0.0/16")},
					},
				},
			},
			RoutingWeight:                  to.Int32Ptr(0),
			SharedKey:                      to.StringPtr("Abc123"),
			UsePolicyBasedTrafficSelectors: to.BoolPtr(false),
			VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
				Resource: armnetwork.Resource{
					ID:       &gatewayId,
					Location: to.StringPtr(config.Location()),
				},
				Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
					Active: to.BoolPtr(false),
					BgpSettings: &armnetwork.BgpSettings{
						Asn:               to.Int64Ptr(65515),
						BgpPeeringAddress: to.StringPtr("10.0.2.30"),
						PeerWeight:        to.Int32Ptr(0),
					},
					EnableBgp:   to.BoolPtr(false),
					GatewayType: armnetwork.VirtualNetworkGatewayTypeVPN.ToPtr(),
					IPConfigurations: &[]*armnetwork.VirtualNetworkGatewayIPConfiguration{{
						SubResource: armnetwork.SubResource{
							ID: to.StringPtr(gatewayId + "/ipConfigurations/" + ipConfigName + ""),
						},
						Name: &ipConfigName,
						Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
							PublicIPAddress: &armnetwork.SubResource{
								ID: &publicAddressId,
							},
							Subnet: &armnetwork.SubResource{
								ID: &subnetId,
							},
						},
					}},
					SKU: &armnetwork.VirtualNetworkGatewaySKU{
						Name: armnetwork.VirtualNetworkGatewaySKUNameVPNGw1.ToPtr(),
						Tier: armnetwork.VirtualNetworkGatewaySKUTierVPNGw1.ToPtr(),
					},
					VPNType: armnetwork.VPNTypeRouteBased.ToPtr(),
				},
			},
		},
	}

	err = CreateVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayConnectionName, gatewayConnectionParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network gateway connection: % +v", err)
	}
	t.Logf("created virtual network gateway connection")

	err = BeginSetVirtualNetworkGatewayConnectionSharedKey(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to set the virtual network gateway connection shared key: %+v", err)
	}
	t.Logf("set the virtual network gateway connection shared key")

	err = GetVirtualNetworkGatewayConnectionSharedKey(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to get the virtual network gateway connection shared key: %+v", err)
	}
	t.Logf("got the virtual network gateway connection shared key")

	err = GetVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to get virtual network gateway connection: %+v", err)
	}
	t.Logf("got virtual network gateway connection")

	err = ListVirtualNetworkGatewayConnection(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual network gateway connection: %+v", err)
	}
	t.Logf("listed virtual network gateway connection")

	//need to sleep for a period of time to run successfully, otherwise "Another operation on this or dependent resource is in progress." will be reported.
	time.Sleep(time.Duration(60) * time.Second)

	err = BeginResetVirtualNetworkGatewayConnectionSharedKey(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to begin reset the virtual network gateway connection shared key: %+v", err)
	}
	t.Logf("began reset the virtual network gateway connection shared key")

	//need to sleep for a period of time to run successfully, otherwise "Another operation on this or dependent resource is in progress." will be reported.
	time.Sleep(time.Duration(60) * time.Second)

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateVirtualNetworkGatewayConnectionTags(ctx, virtualNetworkGatewayConnectionName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for virtual network gateway connection: %+v", err)
	}
	t.Logf("updated virtual network gateway connection tags")

	err = DeleteVirtualNetworkGatewayConnection(ctx, virtualNetworkGatewayConnectionName)
	if err != nil {
		t.Fatalf("failed to delete virtual network gateway connection: %+v", err)
	}
	t.Logf("deleted virtual network gateway connection")
}
