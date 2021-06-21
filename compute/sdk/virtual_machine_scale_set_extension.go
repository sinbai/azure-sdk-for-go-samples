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

func getVirtualMachineScaleSetExtensionsClient() armcompute.VirtualMachineScaleSetExtensionsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetExtensionsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// The operation to create or update an extension.
func CreateVirtualMachineScaleSetExtension(ctx context.Context, vmScaleSetName string, vmssExtensionName string,
	extensionParameters armcompute.VirtualMachineScaleSetExtension) error {
	client := getVirtualMachineScaleSetExtensionsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		vmScaleSetName,
		vmssExtensionName,
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

// The operation to get the extension
func GetVirtualMachineScaleSetExtension(ctx context.Context, vmScaleSetName string, vmssExtensionName string) error {
	client := getVirtualMachineScaleSetExtensionsClient()
	_, err := client.Get(ctx, config.GroupName(), vmScaleSetName, vmssExtensionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets a list of all extensions in a VM scale set.
func ListVirtualMachineScaleSetExtension(ctx context.Context, vmScaleSetName string) error {
	client := getVirtualMachineScaleSetExtensionsClient()
	pager := client.List(config.GroupName(), vmScaleSetName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// The operation to update an extension.
func UpdateVirtualMachineScaleSetExtensionTags(ctx context.Context, vmScaleSetName string, vmssExtensionName string,
	extensionParameters armcompute.VirtualMachineScaleSetExtensionUpdate) error {
	client := getVirtualMachineScaleSetExtensionsClient()
	poller, err := client.BeginUpdate(
		ctx,
		config.GroupName(),
		vmScaleSetName,
		vmssExtensionName,
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

// The operation to delete the extension.
func DeleteVirtualMachineScaleSetExtension(ctx context.Context, vmScaleSetName string, vmssExtensionName string) error {
	client := getVirtualMachineScaleSetExtensionsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), vmScaleSetName, vmssExtensionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
