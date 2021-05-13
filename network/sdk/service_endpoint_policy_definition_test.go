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
)

func TestServiceEndpointPolicyDefinition(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	config.SetLocation("eastus")

	serviceEndpointPolicyName := config.AppendRandomSuffix("serviceendpointpolicy")
	serviceEndpointPolicyDefinitionName := config.AppendRandomSuffix("serviceendpointpolicydefinition")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.DefaultLocation())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateServiceEndpointPolicy(ctx, serviceEndpointPolicyName)
	if err != nil {
		t.Fatalf("failed to create service endpoint policy: % +v", err)
	}
	t.Logf("created service endpoint policy")

	err = CreateServiceEndpointPolicyDefinition(ctx, serviceEndpointPolicyName, serviceEndpointPolicyDefinitionName)
	if err != nil {
		t.Fatalf("failed to create service endpoint policy definition: % +v", err)
	}
	t.Logf("created service endpoint policy definition")

	err = GetServiceEndpointPolicyDefinition(ctx, serviceEndpointPolicyName, serviceEndpointPolicyDefinitionName)
	if err != nil {
		t.Fatalf("failed to get service endpoint policy definition: %+v", err)
	}
	t.Logf("got service endpoint policy definition")

	err = ListServiceEndpointPolicyDefinitionByResourceGroup(ctx, serviceEndpointPolicyName)
	if err != nil {
		t.Fatalf("failed to listservice endpoint policy definitionby resource group: %+v", err)
	}
	t.Logf("listedservice endpoint policy definitionby resource group")

	err = DeleteServiceEndpointPolicyDefinition(ctx, serviceEndpointPolicyName, serviceEndpointPolicyDefinitionName)
	if err != nil {
		t.Fatalf("failed to delete service endpoint policy definition: %+v", err)
	}
	t.Logf("deleted service endpoint policy definition")
}
