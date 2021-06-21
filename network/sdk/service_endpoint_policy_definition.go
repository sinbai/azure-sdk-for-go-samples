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
)

func getServiceEndpointPolicyDefinitionsClient() armnetwork.ServiceEndpointPolicyDefinitionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewServiceEndpointPolicyDefinitionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create ServiceEndpointPolicyDefinitions
func CreateServiceEndpointPolicyDefinition(ctx context.Context, serviceEndpointPolicyName string, serviceEndpointPolicyDefinitionName string, serviceEndpointPolicyDefinitionParameters armnetwork.ServiceEndpointPolicyDefinition) error {
	client := getServiceEndpointPolicyDefinitionsClient()

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		serviceEndpointPolicyName,
		serviceEndpointPolicyDefinitionName,
		serviceEndpointPolicyDefinitionParameters,
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

// Gets the specified service endpoint policy definition in a specified resource group.
func GetServiceEndpointPolicyDefinition(ctx context.Context, serviceEndpointPolicyName string, serviceEndpointPolicyDefinitionName string) error {
	client := getServiceEndpointPolicyDefinitionsClient()
	_, err := client.Get(ctx, config.GroupName(), serviceEndpointPolicyName, serviceEndpointPolicyDefinitionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified service endpoint policy definition.
func DeleteServiceEndpointPolicyDefinition(ctx context.Context, serviceEndpointPolicyName string, serviceEndpointPolicyDefinitionName string) error {
	client := getServiceEndpointPolicyDefinitionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), serviceEndpointPolicyName, serviceEndpointPolicyDefinitionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all service endpoint policy definition in a resource group.
func ListServiceEndpointPolicyDefinitionByResourceGroup(ctx context.Context, serviceEndpointPolicyName string) error {
	client := getServiceEndpointPolicyDefinitionsClient()
	pager := client.ListByResourceGroup(config.GroupName(), serviceEndpointPolicyName, nil)
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
