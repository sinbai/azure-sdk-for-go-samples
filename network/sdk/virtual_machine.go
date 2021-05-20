// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getVirtualMachinesClient() armcompute.VirtualMachinesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualMachines
func CreateVirtualMachine(ctx context.Context, virtualMachineName string, virtualMachinePro armcompute.VirtualMachine) (string, error) {
	client := getVirtualMachinesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualMachineName,
		virtualMachinePro,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.VirtualMachine.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.VirtualMachine.ID, nil
}

// Deletes the specified virtual machine.
func DeleteVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
