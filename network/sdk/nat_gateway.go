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
func CreateNatGateway(ctx context.Context, natGatewayName string, natGatewayParameters armnetwork.NatGateway) error {
	client := getNatGatewayClient()

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		natGatewayName,
		natGatewayParameters,
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
