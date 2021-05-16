// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getVirtualNetworkGatewayConnectionsClient() armnetwork.VirtualNetworkGatewayConnectionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworkGatewayConnectionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualNetworkGatewayConnections
func CreateVirtualNetworkGatewayConnection(ctx context.Context, virtualNetworkName string, virtualNetworkGatewayConnectionName string,
	virtualNetworkGatewayName string, localVirtualNetworkGatewayName string, pipAddressName string, gatewaySubNetName string, ipConfigName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualNetworkGatewayConnectionName,
		armnetwork.VirtualNetworkGatewayConnection{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
				ConnectionProtocol: armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2.ToPtr(),
				ConnectionType:     armnetwork.VirtualNetworkGatewayConnectionTypeIPsec.ToPtr(),
				EnableBgp:          to.BoolPtr(false),
				IPSecPolicies:      &[]*armnetwork.IPSecPolicy{},
				LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
					Resource: armnetwork.Resource{
						ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/localNetworkGateways/" + localVirtualNetworkGatewayName + ""),
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
				TrafficSelectorPolicies:        &[]*armnetwork.TrafficSelectorPolicy{},
				UsePolicyBasedTrafficSelectors: to.BoolPtr(false),
				VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
					Resource: armnetwork.Resource{
						ID:       to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/virtualNetworkGateways/" + virtualNetworkGatewayName + ""),
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
								ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/virtualNetworkGateways/" + virtualNetworkGatewayName + "/ipConfigurations/" + ipConfigName + ""),
							},
							Name: &ipConfigName,
							Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
								PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
								PublicIPAddress: &armnetwork.SubResource{
									ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/publicIPAddresses/" + pipAddressName + ""),
								},
								Subnet: &armnetwork.SubResource{
									ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/virtualNetworks/" + virtualNetworkName + "/subnets/" + gatewaySubNetName),
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

// BeginSetSharedKey - The Put VirtualNetworkGatewayConnectionSharedKey operation sets the virtual network gateway connection shared key for passed virtual
// network gateway connection in the specified resource group through
// Network resource provider.
func BeginSetVirtualNetworkGatewayConnectionSharedKey(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.BeginSetSharedKey(ctx, config.GroupName(), virtualNetworkGatewayConnectionName,
		armnetwork.ConnectionSharedKey{Value: to.StringPtr("AzureAbc124")}, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetSharedKey - The Get VirtualNetworkGatewayConnectionSharedKey operation retrieves information about the specified virtual network gateway connection
// shared key through Network resource provider.
func GetVirtualNetworkGatewayConnectionSharedKey(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.GetSharedKey(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified virtual network gateway connection in a specified resource group.
func GetVirtualNetworkGatewayConnection(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.Get(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the virtual network gateway connection in a subscription.
func ListVirtualNetworkGatewayConnection(ctx context.Context) error {
	client := getVirtualNetworkGatewayConnectionsClient()
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

// BeginResetSharedKey - The VirtualNetworkGatewayConnectionResetSharedKey operation resets the virtual network gateway connection shared key for passed
// virtual network gateway connection in the specified resource group
// through Network resource provider.
func BeginResetVirtualNetworkGatewayConnectionSharedKey(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.BeginResetSharedKey(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, armnetwork.ConnectionResetSharedKey{
		KeyLength: to.Int32Ptr(128),
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Updates virtual network gateway connection tags.
func UpdateVirtualNetworkGatewayConnectionTags(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	_, err := client.BeginUpdateTags(
		ctx,
		config.GroupName(),
		virtualNetworkGatewayConnectionName,
		armnetwork.TagsObject{
			Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
		},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified virtual network gateway connection.
func DeleteVirtualNetworkGatewayConnection(ctx context.Context, virtualNetworkGatewayConnectionName string) error {
	client := getVirtualNetworkGatewayConnectionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualNetworkGatewayConnectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
