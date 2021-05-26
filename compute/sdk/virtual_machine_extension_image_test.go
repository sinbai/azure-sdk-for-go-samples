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

func TestVirtualMachineExtensionImage(t *testing.T) {
	groupName := config.GenerateGroupName("compute")
	config.SetGroupName(groupName)

	extensionPublisherName := "Microsoft.Compute"
	extensionImageType := "VMAccessAgent"
	extensionImageVersion := "1.0.2"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = GetVirtualMachineExtensionImage(ctx, extensionPublisherName, extensionImageType, extensionImageVersion)
	if err != nil {
		t.Fatalf("failed to get virtual machine extension image: %+v", err)
	}
	t.Logf("got virtual machine extension image")

	err = ListVirtualMachineExtensionImageType(ctx, extensionPublisherName)
	if err != nil {
		t.Fatalf("failed to list virtual machine extension image type: %+v", err)
	}
	t.Logf("listed virtual machine extension image type")

	err = ListVirtualMachineExtensionImageVersion(ctx, extensionPublisherName, extensionImageType)
	if err != nil {
		t.Fatalf("failed to list virtual machine extension image version: %+v", err)
	}
	t.Logf("listed virtual machine extension image version")
}
