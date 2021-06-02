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
)

func getSecurityPartnerProvidersClient() armnetwork.SecurityPartnerProvidersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSecurityPartnerProvidersClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates the specified Security Partner Provider
func CreateSecurityPartnerProvider(ctx context.Context, securityPartnerProviderName string, securityPartnerProviderParameters armnetwork.SecurityPartnerProvider) error {
	client := getSecurityPartnerProvidersClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		securityPartnerProviderName,
		securityPartnerProviderParameters,
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

// Gets the specified Security Partner Provider.
func GetSecurityPartnerProvider(ctx context.Context, securityPartnerProviderName string) error {
	client := getSecurityPartnerProvidersClient()
	_, err := client.Get(ctx, config.GroupName(), securityPartnerProviderName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the Security Partner Providers in a subscription.
func ListSecurityPartnerProvider(ctx context.Context) error {
	client := getSecurityPartnerProvidersClient()
	pager := client.List(nil)

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

// Updates tags of a Security Partner Provider resource
func UpdateSecurityPartnerProviderTags(ctx context.Context, securityPartnerProviderName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getSecurityPartnerProvidersClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		securityPartnerProviderName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

//  Deletes the specified Security Partner Provider
func DeleteSecurityPartnerProvider(ctx context.Context, securityPartnerProviderName string) error {
	client := getSecurityPartnerProvidersClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), securityPartnerProviderName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Lists all Security Partner Providers in a resource group.
func ListSecurityPartnerProviderByResourceGroup(ctx context.Context) error {
	client := getSecurityPartnerProvidersClient()
	pager := client.ListByResourceGroup(config.GroupName(), nil)
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
