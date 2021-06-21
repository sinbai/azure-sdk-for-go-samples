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

func getVirtualMachineScaleSetRollingUpgradesClient() armcompute.VirtualMachineScaleSetRollingUpgradesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetRollingUpgradesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Starts a rolling upgrade to move all extensions for all virtual machine scale set instances to the latest available extension
// version. Instances which are already running the latest extension versions
// are not affected.
func StartRollingExtensionUpgrade(ctx context.Context, vmScaleSetName string) error {
	client := getVirtualMachineScaleSetRollingUpgradesClient()
	poller, err := client.BeginStartExtensionUpgrade(ctx, config.GroupName(), vmScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Starts a rolling upgrade to move all virtual machine scale set instances to the latest available Platform Image OS version. Instances
// which are already running the latest available OS version are not
// affected.
func StartRollingUpgradeOSUpgrade(ctx context.Context, vmScaleSetName string, polluntilDown bool) error {
	client := getVirtualMachineScaleSetRollingUpgradesClient()
	poller, err := client.BeginStartOSUpgrade(ctx, config.GroupName(), vmScaleSetName, nil)

	if err != nil {
		return err
	}

	// do not call PollUntilDone function for cancel test
	if polluntilDown {
		_, err = poller.PollUntilDone(ctx, 30*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

// Gets the status of the latest virtual machine scale set rolling upgrade.
func GetLatestVirtualMachineScaleSetRollingUpgrade(ctx context.Context, vmScaleSetName string) error {
	client := getVirtualMachineScaleSetRollingUpgradesClient()
	_, err := client.GetLatest(ctx, config.GroupName(), vmScaleSetName, nil)

	if err != nil {
		return err
	}
	return nil
}

// Cancels the current virtual machine scale set rolling upgrade
func CancelScaleSetRollingUpgrade(ctx context.Context, vmScaleSetName string) error {
	client := getVirtualMachineScaleSetRollingUpgradesClient()
	poller, err := client.BeginCancel(ctx, config.GroupName(), vmScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
