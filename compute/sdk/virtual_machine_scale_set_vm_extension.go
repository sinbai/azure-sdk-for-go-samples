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

func getVirtualMachineScaleSetVmExtensionsClient() armcompute.VirtualMachineExtensionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineExtensionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// The operation to create or update the extension
func CreateVirtualMachineScaleSetVmExtension(ctx context.Context, vmName string, vmExtensionName string, extensionParameters armcompute.VirtualMachineExtension) error {
	client := getVirtualMachineScaleSetVmExtensionsClient()
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

// get the extension
func GetVirtualMachineScaleSetVmExtension(ctx context.Context, vmName string, vmExtensionName string) error {
	client := getVirtualMachineScaleSetVmExtensionsClient()
	_, err := client.Get(ctx, config.GroupName(), vmName, vmExtensionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// The operation to get all extensions of a Virtual Machine
func ListVirtualMachineScaleSetVmExtension(ctx context.Context, vmName string) error {
	client := getVirtualMachineScaleSetVmExtensionsClient()
	_, err := client.List(ctx, config.GroupName(), vmName, nil)

	if err != nil {
		return err
	}
	return nil
}

// The operation to update the extension.
func UpdateVirtualMachineScaleSetVmExtension(ctx context.Context, vmName string, vmExtensionName string, extensionParameters armcompute.VirtualMachineExtensionUpdate) error {
	client := getVirtualMachineScaleSetVmExtensionsClient()
	_, err := client.BeginUpdate(
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
	return nil
}

//  The operation to delete the extension.
func DeleteVirtualMachineScaleSetVmExtension(ctx context.Context, vmName string, vmExtensionName string) error {
	client := getVirtualMachineScaleSetVmExtensionsClient()
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
