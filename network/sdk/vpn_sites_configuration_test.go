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
	"github.com/Azure/go-autorest/autorest/to"
)

func TestVpnSitesConfiguration(t *testing.T) {
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

	request := armnetwork.GetVPNSitesConfigurationRequest{
		OutputBlobSasURL: to.StringPtr("https://blobcortextesturl.blob.core.windows.net/folderforconfig/vpnFile?sp=rw&se=2018-01-10T03%3A42%3A04Z&sv=2017-04-17&sig=WvXrT5bDmDFfgHs%2Brz%2BjAu123eRCNE9BO0eQYcPDT7pY%3D&sr=b"),
		VPNSites:         []*string{&vpnSiteId},
	}
	err = DownloadVpnSitesConfiguration(ctx, virtualWanName, request)
	if err != nil {
		t.Fatalf("failed to download the configurations for vpn-sites in a resource group: %+v", err)
	}
	t.Logf("downloaded the configurations for vpn-sites in a resource group")

}
