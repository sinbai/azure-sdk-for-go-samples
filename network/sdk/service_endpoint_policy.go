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
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getServiceEndpointPoliciesClient() armnetwork.ServiceEndpointPoliciesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewServiceEndpointPoliciesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

//  Creates or updates a service Endpoint Policies.
func CreateServiceEndpointPolicy(ctx context.Context, serviceEndpointPolicyName string) error {
	client := getServiceEndpointPoliciesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		serviceEndpointPolicyName,
		armnetwork.ServiceEndpointPolicy{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
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

// Gets the specified service endpoint policy in a specified resource group.
func GetServiceEndpointPolicy(ctx context.Context, serviceEndpointPolicyName string) error {
	client := getServiceEndpointPoliciesClient()
	_, err := client.Get(ctx, config.GroupName(), serviceEndpointPolicyName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the service endpoint policies in a subscription.
func ListServiceEndpointPolicy(ctx context.Context) error {
	client := getServiceEndpointPoliciesClient()
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

// Updates service endpoint policy tags.
func UpdateServiceEndpointPolicyTags(ctx context.Context, serviceEndpointPolicyName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getServiceEndpointPoliciesClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		serviceEndpointPolicyName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified service endpoint policy.
func DeleteServiceEndpointPolicy(ctx context.Context, serviceEndpointPolicyName string) error {
	client := getServiceEndpointPoliciesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), serviceEndpointPolicyName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all service endpoint policy in a resource group.
func ListServiceEndpointPolicyByResourceGroup(ctx context.Context) error {
	client := getServiceEndpointPoliciesClient()
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
