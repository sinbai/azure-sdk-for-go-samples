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

func TestPublicIPAddress(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	publicIpAddressName := config.GenerateGroupName("publicipaddress")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreatePublicIPAddress(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}
	t.Logf("created public ip address")

	err = GetPublicIPAddress(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to get public ip address: %+v", err)
	}
	t.Logf("got public ip address")

	err = ListPublicIPAddress(ctx)
	if err != nil {
		t.Fatalf("failed to list public ip address: %+v", err)
	}
	t.Logf("listed public ip address")

	err = ListAllPublicIPAddress(ctx)
	if err != nil {
		t.Fatalf("failed to list all public ip address: %+v", err)
	}
	t.Logf("listed all public ip address")

	err = UpdateAddressTags(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to update public ip address: %+v", err)
	}
	t.Logf("updated address tags")

	err = DeletePublicIPAddress(ctx, publicIpAddressName)
	if err != nil {
		t.Fatalf("failed to delete public ip address: %+v", err)
	}
	t.Logf("deleted public ip address")
}
