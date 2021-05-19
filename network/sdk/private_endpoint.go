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

func getPrivateEndpointsClient() armnetwork.PrivateEndpointsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPrivateEndpointsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create PrivateEndpoints
func CreatePrivateEndpoint(ctx context.Context, privateEndpointName string, serviceName string, virtualNetworkName string, subNetName string) error {
	client := getPrivateEndpointsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		privateEndpointName,
		armnetwork.PrivateEndpoint{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.PrivateEndpointProperties{
				PrivateLinkServiceConnections: &[]*armnetwork.PrivateLinkServiceConnection{{
					Name: &serviceName,
					Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
						PrivateLinkServiceID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/privateLinkServices/" + serviceName),
					},
				}},
				Subnet: &armnetwork.Subnet{
					SubResource: armnetwork.SubResource{
						ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/virtualNetworks/" + virtualNetworkName + "/subnets/" + subNetName),
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

// Gets the specified private endpoint in a specified resource group.
func GetPrivateEndpoint(ctx context.Context, privateEndpointName string) error {
	client := getPrivateEndpointsClient()
	_, err := client.Get(ctx, config.GroupName(), privateEndpointName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the private endpoint in a subscription.
func ListPrivateEndpoint(ctx context.Context) error {
	client := getPrivateEndpointsClient()
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

// Gets all the private endpoint in a subscription.
func ListAllPrivateEndpointBySubscription(ctx context.Context) error {
	client := getPrivateEndpointsClient()
	pager := client.ListBySubscription(nil)
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

// Deletes the specified private endpoint.
func DeletePrivateEndpoint(ctx context.Context, privateEndpointName string) error {
	client := getPrivateEndpointsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), privateEndpointName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
