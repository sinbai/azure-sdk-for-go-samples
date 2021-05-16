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

func TestLocalNetworkGateway(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	localNetworkGatewayName := config.AppendRandomSuffix("localnetworkgateway")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateLocalNetworkGateway(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to create local network gateway: % +v", err)
	}
	t.Logf("created local network gateway")

	err = GetLocalNetworkGateway(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to get local network gateway: %+v", err)
	}
	t.Logf("got local network gateway")
	err = ListLocalNetworkGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list local network gateway: %+v", err)
	}
	t.Logf("listed local network gateway")
	err = UpdateLocalNetworkGatewayTags(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to update tags for local network gateway: %+v", err)
	}
	t.Logf("updated local network gateway tags")

	err = DeleteLocalNetworkGateway(ctx, localNetworkGatewayName)
	if err != nil {
		t.Fatalf("failed to delete local network gateway: %+v", err)
	}
	t.Logf("deleted local network gateway")

}
