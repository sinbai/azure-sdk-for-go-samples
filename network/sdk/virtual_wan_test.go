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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVirtualWan(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	config.SetLocation("West US")
	virtualWanName := config.AppendRandomSuffix("virtualwan")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)
	defer config.SetLocation(config.DefaultLocation())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWANParameters := armnetwork.VirtualWAN{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
			Tags:     &map[string]*string{"key1": to.StringPtr("value1")},
		},
		Properties: &armnetwork.VirtualWanProperties{
			DisableVPNEncryption: to.BoolPtr(false),
			Type:                 to.StringPtr("Basic"),
		},
	}
	_, err = CreateVirtualWan(ctx, virtualWanName, virtualWANParameters)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}
	t.Logf("created virtual wan")
}
