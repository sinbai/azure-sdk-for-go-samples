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

func TestVirtualHub(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	virtualWanName := config.AppendRandomSuffix("virtualwan")
	virtualHubName := config.AppendRandomSuffix("virtualhub")

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualWanID, err := CreateVirtualWan(ctx, virtualWanName)
	if err != nil {
		t.Fatalf("failed to create virtual wan: % +v", err)
	}
	t.Logf("created virtual wan")

	err = CreateVirtualHub(ctx, virtualHubName, virtualWanID)
	if err != nil {
		t.Fatalf("failed to create virtual hub: % +v", err)
	}
	t.Logf("created virtual hub")

}
