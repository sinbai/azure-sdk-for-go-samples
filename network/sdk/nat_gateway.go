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
	"github.com/Azure/go-autorest/autorest/to"
)

func getNatGatewayClient() armnetwork.NatGatewaysClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNatGatewaysClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a nat gateway.
func CreateNatGateway(ctx context.Context, natGatewayName string, pipaddress string, pipprefix string) error {
	client := getNatGatewayClient()

	urlPathAddress := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{pipaddress}"
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{pipaddress}", url.PathEscape(pipaddress))
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	urlPathPrefix := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/PublicIPPrefixes/{pipprefix}"
	urlPathPrefix = strings.ReplaceAll(urlPathPrefix, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathPrefix = strings.ReplaceAll(urlPathPrefix, "{pipprefix}", url.PathEscape(pipprefix))
	urlPathPrefix = strings.ReplaceAll(urlPathPrefix, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		natGatewayName,
		armnetwork.NatGateway{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.NatGatewayPropertiesFormat{
				PublicIPAddresses: &[]*armnetwork.SubResource{
					{
						ID: &urlPathAddress,
					},
				},
				PublicIPPrefixes: &[]*armnetwork.SubResource{
					{
						ID: &urlPathPrefix,
					},
				},
			},
			SKU: &armnetwork.NatGatewaySKU{
				Name: armnetwork.NatGatewaySKUNameStandard.ToPtr(),
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

// Gets the specified nat gateway in a specified resource group.
func GetNatGateway(ctx context.Context, natGatewayName string) error {
	client := getNatGatewayClient()
	_, err := client.Get(ctx, config.GroupName(), natGatewayName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all nat gateways in a resource group.
func ListNatGateway(ctx context.Context) error {
	client := getNatGatewayClient()
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

// Gets all the Nat Gateways in a subscription.
func ListAllNatGateway(ctx context.Context) error {
	client := getNatGatewayClient()
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

// Updates nat gateway tags.
func UpdateNatGateway(ctx context.Context, natGatewayName string) error {
	client := getNatGatewayClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		natGatewayName,
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

// Deletes the specified nat gateway.
func DeleteNatGatewayGroup(ctx context.Context, natGatewayName string) error {
	client := getNatGatewayClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), natGatewayName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
