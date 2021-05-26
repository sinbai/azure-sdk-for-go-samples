// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

func TestVirtualMachineImage(t *testing.T) {
	groupName := config.GenerateGroupName("compute")
	config.SetGroupName(groupName)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	publisherName := "MicrosoftWindowsServer"
	offer := "WindowsServer"
	skus := "2019-Datacenter"
	version := "2019.0.20190115"

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = GetVirtualMachineImage(ctx, publisherName, offer, skus, version)
	if err != nil {
		t.Fatalf("failed to get virtual machine image: %+v", err)
	}
	t.Logf("got virtual machine image")

	err = ListVirtualMachineImage(ctx, publisherName, offer, skus)
	if err != nil {
		t.Fatalf("failed to list virtual machine image: %+v", err)
	}
	t.Logf("listed virtual machine image")

	err = ListVirtualMachineImageOffer(ctx, publisherName)
	if err != nil {
		t.Fatalf("failed to list virtual machine image offer: %+v", err)
	}
	t.Logf("listed virtual machine image offer")

	err = LisVirtualMachineImagePublisher(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual machine image publisher: %+v", err)
	}
	t.Logf("listed virtual machine image publisher")

	err = ListVirtualMachineImageSKU(ctx, publisherName, offer)
	if err != nil {
		t.Fatalf("failed to list virtual machine image SKU: %+v", err)
	}
	t.Logf("listed virtual machine image SKU")

}
