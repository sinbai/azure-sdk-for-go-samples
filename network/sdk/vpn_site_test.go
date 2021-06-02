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

func TestVpnSite(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	vpnSiteName := config.AppendRandomSuffix("vpnsite")
	virtualWanName := config.AppendRandomSuffix("virtualwan")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
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
	_, err = CreateVpnSite(ctx, vpnSiteName, vpnSiteParameters)
	if err != nil {
		t.Fatalf("failed to create vpn site: % +v", err)
	}
	t.Logf("created vpn site")

	err = ListVpnSite(ctx)
	if err != nil {
		t.Fatalf("failed to list vpn site: %+v", err)
	}
	t.Logf("listed vpn site")

	err = ListVpnSiteByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to listvpn site by resource group: %+v", err)
	}
	t.Logf("listedvpn site by resource group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateVpnSiteTags(ctx, vpnSiteName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for vpn site: %+v", err)
	}
	t.Logf("updated vpn site tags")

	err = DeleteVpnSite(ctx, vpnSiteName)
	if err != nil {
		t.Fatalf("failed to delete vpn site: %+v", err)
	}
	t.Logf("deleted vpn site")
}
