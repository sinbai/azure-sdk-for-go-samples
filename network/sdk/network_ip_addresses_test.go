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

func TestNetwork(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	publicIpPrefixName := config.GenerateGroupName("publicipprefix")
	publicIpAddressName := config.GenerateGroupName("publicipaddress")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.DeleteGroup(ctx, config.GroupName())

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}
	t.Logf("created group: %s\n", groupName)

	CreatePublicIPPrefix(ctx, publicIpPrefixName)
	t.Logf("created public ip prefix")

	CreatePublicIPAddress(ctx, publicIpAddressName)
	t.Logf("created public ip address")

	GetPublicIPPrefix(ctx, publicIpPrefixName)
	t.Logf("got public ip prefix")

	GetPublicIPAddress(ctx, publicIpAddressName)
	t.Logf("got public ip address")

	ListPublicIPPrefix(ctx)
	t.Logf("listed public ip prefix")

	ListPublicIPAddress(ctx)
	t.Logf("listed public ip address")

	ListAllPublicIPPrefix(ctx)
	t.Logf("listed all public ip prefix")

	ListAllPublicIPAddress(ctx)
	t.Logf("listed all public ip address\n")

	UpdatePrefixTags(ctx, publicIpPrefixName)
	t.Logf("updated prefix tags")

	UpdateAddressTags(ctx, publicIpAddressName)
	t.Logf("updated address tags")

	DeletePublicIPPrefix(ctx, publicIpPrefixName)
	t.Logf("deleted public ip prefix")

	DeletePublicIPAddress(ctx, publicIpAddressName)
	t.Logf("deleted public ip address")
}
