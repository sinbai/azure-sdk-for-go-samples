// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
)

func getVirtualMachineExtensionsClient() armcompute.VirtualMachineExtensionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineExtensionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualMachineExtensions
func CreateVirtualMachineExtension(ctx context.Context, vmName string, vmExtensionName string, extensionParameters armcompute.VirtualMachineExtension) error {
	client := getVirtualMachineExtensionsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		vmName,
		vmExtensionName,
		extensionParameters,
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified virtual machine extension in a specified resource group.
func GetVirtualMachineExtension(ctx context.Context, vmName string, vmExtensionName string) error {
	client := getVirtualMachineExtensionsClient()
	_, err := client.Get(ctx, config.GroupName(), vmName, vmExtensionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the virtual machine extension in a subscription.
func ListVirtualMachineExtension(ctx context.Context, vmName string) error {
	client := getVirtualMachineExtensionsClient()
	_, err := client.List(ctx, config.GroupName(), vmName, nil)

	if err != nil {
		return err
	}
	return nil
}

// Updates virtual machine extension tags.
func UpdateVirtualMachineExtensionTags(ctx context.Context, vmName string, vmExtensionName string, extensionParameters armcompute.VirtualMachineExtensionUpdate) error {
	client := getVirtualMachineExtensionsClient()
	poller, err := client.BeginUpdate(
		ctx,
		config.GroupName(),
		vmName,
		vmExtensionName,
		extensionParameters,
		nil,
	)
	if err != nil {
		return err
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified virtual machine extension.
func DeleteVirtualMachineExtension(ctx context.Context, vmName string, vmExtensionName string) error {
	client := getVirtualMachineExtensionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), vmName, vmExtensionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
