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

func TestIPGroup(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	ipGroupName := config.AppendRandomSuffix("ipgroup")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateIPGroup(ctx, ipGroupName)
	if err != nil {
		t.Fatalf("failed to create ip group: %+v", err)
	}
	t.Logf("created ip group")

	err = GetIPGroup(ctx, ipGroupName)
	if err != nil {
		t.Fatalf("failed to get ip group: %+v", err)
	}
	t.Logf("got ip group")

	err = ListIPGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list ip group: %+v", err)
	}
	t.Logf("listed ip group")

	err = ListIPGroupByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to list ip group by resource group: %+v", err)
	}
	t.Logf("listed ip group by resource group")

	err = DeletePublicIPAddress(ctx, ipGroupName)
	if err != nil {
		t.Fatalf("failed to delete ip group: %+v", err)
	}
	t.Logf("deleted ip group")
}
