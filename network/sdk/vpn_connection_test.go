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
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVpnConnection(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	vpnConnectionName := config.AppendRandomSuffix("vpnconnection")
	vpnGatewayName := config.AppendRandomSuffix("vpngateway")
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	vpnSiteName := config.AppendRandomSuffix("vpnsite")
	vpnSiteLinkName := "vpnSiteLink1"

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWANParameters := armnetwork.VirtualWAN{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualWanProperties{
			DisableVPNEncryption: to.BoolPtr(false),
			Type:                 to.StringPtr("Basic"),
		},
	}
	virtualWanId, err := CreateVirtualWan(ctx, virtualWanName, virtualWANParameters)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}

	vpnSiteParameters := armnetwork.VPNSite{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VPNSiteProperties{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{to.StringPtr("10.0.0.0/16")},
			},
			IsSecuritySite: to.BoolPtr(false),
			VPNSiteLinks: []*armnetwork.VPNSiteLink{{
				Name: to.StringPtr("vpnSiteLink1"),
				Properties: &armnetwork.VPNSiteLinkProperties{
					IPAddress: to.StringPtr("50.50.50.56"),
					LinkProperties: &armnetwork.VPNLinkProviderProperties{
						LinkProviderName: to.StringPtr("vendor1"),
						LinkSpeedInMbps:  to.Int32Ptr(0),
					},
					BgpProperties: &armnetwork.VPNLinkBgpSettings{
						Asn:               to.Int64Ptr(1234),
						BgpPeeringAddress: to.StringPtr("192.168.0.0"),
					},
				},
			}},
			VirtualWan: &armnetwork.SubResource{
				ID: &virtualWanId,
			},
		},
	}
	vpnSiteId, err := CreateVpnSite(ctx, vpnSiteName, vpnSiteParameters)
	if err != nil {
		t.Fatalf("failed to create vpn site: % +v", err)
	}

	virtualHubParameters := armnetwork.VirtualHub{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualHubProperties{
			AddressPrefix: to.StringPtr("10.168.0.0/24"),
			SKU:           to.StringPtr("Basic"),
			VirtualWan: &armnetwork.SubResource{
				ID: &virtualWanId,
			},
		},
	}

	virtualHubId, err := CreateVirtualHub(ctx, virtualHubName, virtualHubParameters, false)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}

	vpnGatewayParameters := armnetwork.VPNGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VPNGatewayProperties{
			BgpSettings: &armnetwork.BgpSettings{
				Asn:        to.Int64Ptr(65515),
				PeerWeight: to.Int32Ptr(0),
			},
			Connections: []*armnetwork.VPNConnection{
				{
					Name: &vpnConnectionName,
					Properties: &armnetwork.VPNConnectionProperties{
						RemoteVPNSite: &armnetwork.SubResource{
							ID: &vpnSiteId,
						},
						VPNLinkConnections: []*armnetwork.VPNSiteLinkConnection{
							{
								Name: to.StringPtr("Connection-Link1"),
								Properties: &armnetwork.VPNSiteLinkConnectionProperties{
									ConnectionBandwidth:       to.Int32Ptr(200),
									SharedKey:                 to.StringPtr("key"),
									VPNConnectionProtocolType: armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2.ToPtr(),
									VPNSiteLink: &armnetwork.SubResource{
										ID: to.StringPtr(vpnSiteId + "/vpnSiteLinks/" + vpnSiteLinkName),
									},
								},
							},
						},
					},
				},
			},
			VirtualHub: &armnetwork.SubResource{
				ID: &virtualHubId,
			},
		},
	}
	err = CreateVpnGateway(ctx, vpnGatewayName, vpnGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create vpn gateway: % +v", err)
	}

	vpnConnectionParameters := armnetwork.VPNConnection{
		Properties: &armnetwork.VPNConnectionProperties{
			RemoteVPNSite: &armnetwork.SubResource{
				ID: &vpnSiteId,
			},
			VPNLinkConnections: []*armnetwork.VPNSiteLinkConnection{{
				Name: &vpnConnectionName,
				Properties: &armnetwork.VPNSiteLinkConnectionProperties{
					ConnectionBandwidth:            to.Int32Ptr(200),
					UsePolicyBasedTrafficSelectors: to.BoolPtr(false),
					SharedKey:                      to.StringPtr("key"),
					VPNConnectionProtocolType:      armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2.ToPtr(),
					VPNSiteLink: &armnetwork.SubResource{
						ID: to.StringPtr(vpnSiteId + "/vpnSiteLinks/" + vpnSiteLinkName),
					},
				},
			}},
		},
	}
	err = CreateVpnConnection(ctx, vpnGatewayName, vpnConnectionName, vpnConnectionParameters)
	if err != nil {
		t.Fatalf("failed to create vpn connection: % +v", err)
	}
	t.Logf("created vpn connection")

	err = GetVpnConnection(ctx, vpnGatewayName, vpnConnectionName)
	if err != nil {
		t.Fatalf("failed to get vpn connection: %+v", err)
	}
	t.Logf("got vpn connection")

	err = ListVpnConnectionByVpnGateway(ctx, vpnGatewayName)
	if err != nil {
		t.Fatalf("failed to list vpn connection by resource group: %+v", err)
	}
	t.Logf("listed vpn connection by resource group")

	err = DeleteVpnConnection(ctx, vpnGatewayName, vpnConnectionName)
	if err != nil {
		t.Fatalf("failed to delete vpn connection: %+v", err)
	}
	t.Logf("deleted vpn connection")
}
