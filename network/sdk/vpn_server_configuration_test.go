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

func TestVpnServerConfiguration(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	vpnServerConfigurationName := config.AppendRandomSuffix("vpnserverconfiguration")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
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
			VPNAuthenticationTypes: []*armnetwork.VPNAuthenticationType{armnetwork.VPNAuthenticationTypeAAD.ToPtr()},
		},
	}
	_, err = CreateVpnServerConfiguration(ctx, vpnServerConfigurationName, vpnServerConfigurationParameters)
	if err != nil {
		t.Fatalf("failed to create vpn server configuration: % +v", err)
	}
	t.Logf("created vpn server configuration")

	err = GetVpnServerConfiguration(ctx, vpnServerConfigurationName)
	if err != nil {
		t.Fatalf("failed to get vpn server configuration: %+v", err)
	}
	t.Logf("got vpn server configuration")

	err = ListVpnServerConfiguration(ctx)
	if err != nil {
		t.Fatalf("failed to list vpn server configuration: %+v", err)
	}
	t.Logf("listed vpn server configuration")

	err = ListVpnServerConfigurationByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to listvpn server configuration by resource group: %+v", err)
	}
	t.Logf("listedvpn server configuration by resource group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateVpnServerConfigurationTags(ctx, vpnServerConfigurationName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for vpn server configuration: %+v", err)
	}
	t.Logf("updated vpn server configuration tags")

	err = DeleteVpnServerConfiguration(ctx, vpnServerConfigurationName)
	if err != nil {
		t.Fatalf("failed to delete vpn server configuration: %+v", err)
	}
	t.Logf("deleted vpn server configuration")

}
