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

func TestPublicIPPrefix(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	publicIpPrefixName := config.AppendRandomSuffix("pipperfix")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreatePublicIPPrefix(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to create public ip prefix: %+v", err)
	}
	t.Logf("created public ip prefix")

	err = GetPublicIPPrefix(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to get public ip prefix: %+v", err)
	}
	t.Logf("got public ip prefix")

	err = ListPublicIPPrefix(ctx)
	if err != nil {
		t.Fatalf("failed to list public ip prefix: %+v", err)
	}
	t.Logf("listed public ip prefix")

	err = ListAllPublicIPPrefix(ctx)
	if err != nil {
		t.Fatalf("failed to list all public ip prefix: %+v", err)
	}
	t.Logf("listed all public ip prefix")

	err = UpdatePublicIPPrefixTags(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to update tags for public ip prefix: %+v", err)
	}
	t.Logf("updated public ip prefix tags")

	err = DeletePublicIPPrefix(ctx, publicIpPrefixName)
	if err != nil {
		t.Fatalf("failed to delete public ip prefix: %+v", err)
	}
	t.Logf("deleted public ip prefix")
}
