// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getVirtualWansClient() armnetwork.VirtualWansClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualWansClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualWans
func CreateVirtualWan(ctx context.Context, virtualWanName string) (string, error) {
	client := getVirtualWansClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualWanName,
		armnetwork.VirtualWan{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
				Tags:     &map[string]string{"key1": "value1"},
			},
			Properties: &armnetwork.VirtualWanProperties{
				DisableVpnEncryption: to.BoolPtr(false),
				Type:                 to.StringPtr("Basic"),
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}
	return "", nil
}

func getVirtualHubsClient() armnetwork.VirtualHubsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualHubsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualHubs
func CreateVirtualHub(ctx context.Context, virtualHubName string, virtualWanName string) error {
	client := getVirtualHubsClient()

	urlPathVirtualWan := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{virtualWanName}"
	urlPathVirtualWan = strings.ReplaceAll(urlPathVirtualWan, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathVirtualWan = strings.ReplaceAll(urlPathVirtualWan, "{virtualWanName}", url.PathEscape(virtualWanName))
	urlPathVirtualWan = strings.ReplaceAll(urlPathVirtualWan, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualHubName,
		armnetwork.VirtualHub{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
				Tags:     &map[string]string{"key1": "value1"},
			},
			Properties: &armnetwork.VirtualHubProperties{
				AddressPrefix: to.StringPtr("10.168.0.0/24"),
				SKU:           to.StringPtr("Basic"),
				VirtualWan: &armnetwork.SubResource{
					ID: &urlPathVirtualWan,
				},
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func getFirewallsClient() armnetwork.AzureFirewallsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewAzureFirewallsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create Firewalls
func CreateFirewall(ctx context.Context, firewallName string, firewallPolicyName string, virtualHubName string) error {
	client := getFirewallsClient()

	urlPathVirtualHub := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualHubs/{virtualHubName}"
	urlPathVirtualHub = strings.ReplaceAll(urlPathVirtualHub, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathVirtualHub = strings.ReplaceAll(urlPathVirtualHub, "{virtualHubName}", url.PathEscape(virtualHubName))
	urlPathVirtualHub = strings.ReplaceAll(urlPathVirtualHub, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	urlPathFirewallPolicy := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/firewallPolicies/{firewallPolicyName}"
	urlPathFirewallPolicy = strings.ReplaceAll(urlPathFirewallPolicy, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathFirewallPolicy = strings.ReplaceAll(urlPathFirewallPolicy, "{firewallPolicyName}", url.PathEscape(firewallPolicyName))
	urlPathFirewallPolicy = strings.ReplaceAll(urlPathFirewallPolicy, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		firewallName,
		armnetwork.AzureFirewall{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
				Tags:     &map[string]string{"key1": "value1"},
			},
			Properties: &armnetwork.AzureFirewallPropertiesFormat{
				SKU: &armnetwork.AzureFirewallSKU{
					Name: armnetwork.AzureFirewallSKUNameAzfwHub.ToPtr(),
					Tier: armnetwork.AzureFirewallSKUTierStandard.ToPtr(),
				},
				VirtualHub: &armnetwork.SubResource{
					ID: &urlPathVirtualHub,
				},
				FirewallPolicy: &armnetwork.SubResource{
					ID: &urlPathFirewallPolicy,
				},
				HubIPAddresses: &armnetwork.HubIPAddresses{
					PublicIPs: &armnetwork.HubPublicIPAddresses{
						Addresses: &[]armnetwork.AzureFirewallPublicIPAddress{},
						Count:     to.Int32Ptr(1),
					},
				},
			},
			Zones: &[]string{},
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 120*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified firewall in a specified resource group.
func GetFirewall(ctx context.Context, firewallName string) error {
	client := getFirewallsClient()
	_, err := client.Get(ctx, config.GroupName(), firewallName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the firewall in a subscription.
func ListFirewall(ctx context.Context) error {
	client := getFirewallsClient()
	pager := client.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Gets all the firewall in a subscription.
func ListAllFirewall(ctx context.Context) error {
	client := getFirewallsClient()
	pager := client.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Updates firewall tags.
func UpdateFirewallTags(ctx context.Context, firewallName string) error {
	client := getFirewallsClient()
	poller, err := client.BeginUpdateTags(
		ctx,
		config.GroupName(),
		firewallName,
		armnetwork.TagsObject{
			Tags: &map[string]string{"tag1": "value1", "tag2": "value2"},
		},
		nil,
	)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified firewall.
func DeleteFirewall(ctx context.Context, firewallName string) error {
	client := getFirewallsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), firewallName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func getAzureFirewallFqdnTagsClient() armnetwork.AzureFirewallFqdnTagsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewAzureFirewallFqdnTagsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets all the azure firewall fqdn tag in a subscription.
func ListAllAzureFirewallFqdnTag(ctx context.Context) error {
	client := getAzureFirewallFqdnTagsClient()
	pager := client.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}
