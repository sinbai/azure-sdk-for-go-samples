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

func TestVpnServerConfigurationsAssociatedWithVirtualWan(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	vpnGatewayName := config.AppendRandomSuffix("vpngateway")
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")
	vpnSiteName := config.AppendRandomSuffix("vpnsite")
	vpnSiteLinkName := "vpnSiteLink1"
	vpnServerConfigurationName := config.AppendRandomSuffix("vpnserverconfiguration")
	p2sVpnGatewayName := config.AppendRandomSuffix("p2sVpnGateway")
	p2sConnectionConfigurationName := config.AppendRandomSuffix("p2sconnectionconfiguration")
	ctx, cancel := context.WithTimeout(context.Background(), 4000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWANParameters := armnetwork.VirtualWAN{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualWanProperties{
			DisableVPNEncryption: to.BoolPtr(false),
			Type:                 to.StringPtr("Standard"),
		},
	}
	virtualWanId, err := CreateVirtualWan(ctx, virtualWanName, virtualWANParameters)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}

	vpnSiteParameters := armnetwork.VPNSite{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VPNSiteProperties{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
			IsSecuritySite: to.BoolPtr(false),
			VPNSiteLinks: &[]*armnetwork.VPNSiteLink{{
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
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualHubProperties{
			AddressPrefix: to.StringPtr("10.168.0.0/24"),
			SKU:           to.StringPtr("Basic"),
			VirtualWan: &armnetwork.SubResource{
				ID: &virtualWanId,
			},
		},
	}

	virtualHubId, err := CreateVirtualHub(ctx, virtualHubName, virtualHubParameters)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}

	vpnGatewayParameters := armnetwork.VPNGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VPNGatewayProperties{
			BgpSettings: &armnetwork.BgpSettings{
				Asn:        to.Int64Ptr(65515),
				PeerWeight: to.Int32Ptr(0),
			},
			Connections: &[]*armnetwork.VPNConnection{
				{
					Name: to.StringPtr("vpnConnection1"),
					Properties: &armnetwork.VPNConnectionProperties{
						RemoteVPNSite: &armnetwork.SubResource{
							ID: &vpnSiteId,
						},
						VPNLinkConnections: &[]*armnetwork.VPNSiteLinkConnection{
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

	vpnServerConfigurationParameters := armnetwork.VPNServerConfiguration{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.VPNServerConfigurationProperties{
			AADAuthenticationParameters: &armnetwork.AADAuthenticationParameters{
				AADAudience: to.StringPtr("00000000-abcd-abcd-abcd-999999999999"),
				AADIssuer:   to.StringPtr("https://sts.windows.net/" + config.TenantID() + "/"),
				AADTenant:   to.StringPtr("https://login.microsoftonline.com/" + config.TenantID()),
			},
			VPNAuthenticationTypes: &[]*armnetwork.VPNAuthenticationType{armnetwork.VPNAuthenticationTypeAAD.ToPtr()},
		},
	}
	vpnServerConfigurationId, err := CreateVpnServerConfiguration(ctx, vpnServerConfigurationName, vpnServerConfigurationParameters)
	if err != nil {
		t.Fatalf("failed to create vpn server configuration: % +v", err)
	}

	p2SVPNGatewayParameters := armnetwork.P2SVPNGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.P2SVPNGatewayProperties{
			P2SConnectionConfigurations: &[]*armnetwork.P2SConnectionConfiguration{{
				SubResource: armnetwork.SubResource{
					ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/p2sVpnGateways/" + p2sVpnGatewayName + "/p2sConnectionConfigurations/" + p2sConnectionConfigurationName),
				},
				Name: to.StringPtr("P2SConnectionConfig1"),
				Properties: &armnetwork.P2SConnectionConfigurationProperties{
					VPNClientAddressPool: &armnetwork.AddressSpace{
						AddressPrefixes: &[]*string{to.StringPtr("101.3.0.0/16")},
					},
				},
			}},
			VPNGatewayScaleUnit: to.Int32Ptr(1),
			VPNServerConfiguration: &armnetwork.SubResource{
				ID: &vpnServerConfigurationId,
			},
			VirtualHub: &armnetwork.SubResource{
				ID: &virtualHubId,
			},
			CustomDNSServers: &[]*string{
				to.StringPtr("1.1.1.1"),
				to.StringPtr("2.2.2.2"),
			},
		},
	}
	err = CreateP2sVpnGateway(ctx, vpnGatewayName, p2SVPNGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create p2s vpn gateway: % +v", err)
	}

	err = ListVpnServerConfigurationsAssociatedWithVirtualWan(ctx, virtualWanName)
	if err != nil {
		t.Fatalf("failed to list vpn server configurations associated with virtual wan: %+v", err)
	}
	t.Logf("listed vpn server configurations associated with virtual wan")

}
