// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
)

func getVirtualMachineRunCommandsClient() armcompute.VirtualMachineRunCommandsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineRunCommandsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets specific run command for a subscription in a location.
func GetVirtualMachineRunCommand(ctx context.Context, commandId string) error {
	client := getVirtualMachineRunCommandsClient()
	_, err := client.Get(ctx, config.Location(), commandId, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all available run commands for a subscription in a location.
func ListVirtualMachineRunCommand(ctx context.Context) error {
	client := getVirtualMachineRunCommandsClient()
	pager := client.List(config.Location(), nil)

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
