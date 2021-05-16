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

func TestServiceAssociationLink(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName := "Subnet"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateVirtualNetwork(ctx, virtualNetworkName, "10.0.0.0/16")
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	body := `{
		"addressPrefix": "10.0.1.0/24"
		}`
	_, err = CreateSubnet(ctx, virtualNetworkName, subNetName, body)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	err = ListServiceAssociationLink(ctx, virtualNetworkName, subNetName)
	if err != nil {
		t.Fatalf("failed to list service association link: %+v", err)
	}
	t.Logf("listed service association link")
}
